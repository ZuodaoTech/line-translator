package line

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/zuodaotech/line-translator/common/uuid"
	"gopkg.in/square/go-jose.v2"
)

type (
	Client struct {
		cfg         Config
		bot         *messaging_api.MessagingApiAPI
		accessToken string
	}
	Config struct {
		ChannelID  string
		ChannelKey string
		PrivateKey string
	}
)

func New(cfg Config) (*Client, error) {
	decoded, err := base64.StdEncoding.DecodeString(cfg.PrivateKey)
	if err != nil {
		log.Fatalf("failed to decode line private key: %v", err)
		return nil, err
	}
	cfg.PrivateKey = string(decoded)

	return &Client{
		cfg: cfg,
		bot: nil,
	}, nil
}

func NewFromAccessToken(token string) (*Client, error) {
	bot, err := messaging_api.NewMessagingApiAPI(token)
	if err != nil {
		return nil, err
	}

	return &Client{
		cfg: Config{},
		bot: bot,
	}, nil
}

func (s *Client) GenerateToken() (string, *time.Time, error) {
	jwt, err := s.GenerateJWTFromJWK(s.cfg.PrivateKey, s.cfg.ChannelKey)
	if err != nil {
		return "", nil, err
	}

	token, expiredAt, err := s.getChannelAccessToken(jwt)
	if err != nil {
		return "", nil, err
	}

	s.bot, err = messaging_api.NewMessagingApiAPI(token)
	if err != nil {
		return "", nil, err
	}

	return token, expiredAt, nil
}

func (s *Client) SendPushMessage(ctx context.Context, groupID, title, summary, url string) error {
	content := fmt.Sprintf("%s\n%s\n\n👉 %s", title, summary, url)
	_, err := s.bot.PushMessage(&messaging_api.PushMessageRequest{
		To: groupID,
		Messages: []messaging_api.MessageInterface{
			messaging_api.TextMessage{
				Text: content,
			},
		},
	}, uuid.New())

	if err != nil {
		return err
	}
	return nil
}

func (s *Client) GenerateJWTFromJWK(jwkJSON string, kid string) (string, error) {
	// Parse the JWK
	var jwk jose.JSONWebKey
	err := jwk.UnmarshalJSON([]byte(jwkJSON))
	if err != nil {
		return "", err
	}

	// Convert JWK to RSA Private Key
	rsaPrivateKey, ok := jwk.Key.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("failed to convert JWK to RSA Private Key")
	}

	// Define the token's claims
	claims := jwt.MapClaims{
		"iss":       s.cfg.ChannelID,
		"sub":       s.cfg.ChannelID,
		"aud":       "https://api.line.me/",
		"exp":       time.Now().Add(time.Minute * 29).Unix(), // 29 minutes from now
		"token_exp": 86400,
	}

	// Create a new token with the specified algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	token.Header["typ"] = "JWT"
	token.Header["alg"] = "RS256"

	// Sign the token with the private key
	return token.SignedString(rsaPrivateKey)
}

func (s *Client) getChannelAccessToken(jwtToken string) (string, *time.Time, error) {
	// LINE API endpoint for obtaining channel access token
	const tokenEndpoint = "https://api.line.me/oauth2/v2.1/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	data.Set("client_assertion", jwtToken)

	// Create request
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", nil, errors.New("failed to get channel access token: " + string(body))
	}

	// Parse response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, err
	}

	// Extract access token
	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", nil, errors.New("access token not found in response")
	}

	expiresInFlt, ok := result["expires_in"].(float64)
	if !ok {
		return "", nil, errors.New("expires_in not found in response")
	}
	expiresIn := int64(expiresInFlt)

	expiredAt := time.Now().Add(time.Duration(expiresIn)*time.Second - 10)

	return accessToken, &expiredAt, nil
}

func (s *Client) ReplyTextMessage(replyToken, quoteToken string, text string) (*messaging_api.ReplyMessageResponse, error) {
	msg := messaging_api.TextMessage{
		Text: text,
	}
	if quoteToken != "" {
		msg.QuoteToken = quoteToken
	}
	resp, err := s.bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				msg,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GenerateJWKPair() (string, string, error) {
	var rawkey interface{}
	attrs := map[string]interface{}{
		"alg": "RS256",
		"use": "sig",
	}
	v, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		slog.Error("[common.line] failed to generate public key", "error", err)
		return "", "", err
	}
	rawkey = v
	key, err := jwk.FromRaw(rawkey)
	if err != nil {
		slog.Error("[common.line] failed to extract public key", "error", err)
		return "", "", err
	}
	for k, v := range attrs {
		if err := key.Set(k, v); err != nil {
			return "", "", err
		}
	}

	keyset := jwk.NewSet()
	keyset.AddKey(key)

	pubks, err := jwk.PublicSetOf(keyset)
	if err != nil {
		return "", "", err
	}

	keybuf, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return "", "", err
	}
	pub, _ := pubks.Key(0)
	pubbuf, err := json.MarshalIndent(pub, "", "  ")
	if err != nil {
		return "", "", err
	}

	// encode to base64
	encodedKey := base64.StdEncoding.EncodeToString(keybuf)
	encodedPub := base64.StdEncoding.EncodeToString(pubbuf)

	return encodedPub, encodedKey, err
}
