package line

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getChannelAccessToken(jwtToken string) (string, *time.Time, error) {
	// LINE API endpoint for obtaining channel access token
	const tokenEndpoint = "https://api.line.me/oauth2/v2.1/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	data.Set("client_assertion", jwtToken)

	return getAccessToken(tokenEndpoint, data)
}

func getChannelStatelessAccessToken(jwtToken string) (string, *time.Time, error) {
	// curl -v -X POST https://api.line.me/oauth2/v3/token \
	// -H 'Content-Type: application/x-www-form-urlencoded' \
	// --data-urlencode 'grant_type=client_credentials' \
	// --data-urlencode 'client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer' \
	// --data-urlencode 'client_assertion={JWT assertion}'

	const tokenEndpoint = "https://api.line.me/oauth2/v3/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	data.Set("client_assertion", jwtToken)

	return getAccessToken(tokenEndpoint, data)
}

func getAccessToken(url string, data url.Values) (string, *time.Time, error) {
	// Create request
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
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
