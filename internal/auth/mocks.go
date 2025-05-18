package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MockAuthenticator struct{}

const secret = "test"

var testClaims = jwt.MapClaims{
	"sub": int64(202),
	"exp": time.Now().Add(time.Hour).Unix(),
	"iss": "test-aud",
	"aud": "test-aud",
}

func (a *MockAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	return token.SignedString([]byte(secret))
}

func (a *MockAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}
