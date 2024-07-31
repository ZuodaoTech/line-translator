package session

import (
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

type Session struct {
	WebBase string

	Version string

	JwtSecret []byte

	rdb *redis.Client
}

type (
	JwtClaims struct {
		Version int    `json:"version"`
		UserID  uint64 `json:"user_id"`
		AppID   string `json:"app_id"`
		Scope   string `json:"scope"`
		jwt.RegisteredClaims
	}
)

func (s *Session) WithJWTSecret(secret []byte) *Session {
	s.JwtSecret = secret
	return s
}

func (s *Session) WithRdb(rdb *redis.Client) *Session {
	s.rdb = rdb
	return s
}
