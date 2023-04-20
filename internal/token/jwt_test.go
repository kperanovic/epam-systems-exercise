package token

import (
	"strings"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = "KLguRWx03zXcWwDXywrxgwTS7r39QaF1"

func TestNewJWTToken(t *testing.T) {
	jwt, err := NewJWTToken(secretKey)
	assert.Equal(t, err, nil)
	assert.Equal(t, jwt, &JWTToken{secretKey: secretKey})

	// Test if key is smaller than 32 characters
	jwt, err = NewJWTToken(strings.Split(secretKey, "W")[0])
	assert.NotEqual(t, err, nil)
	assert.Equal(t, jwt, nil)

}

func TestVerifyToken(t *testing.T) {
	j, err := NewJWTToken(secretKey)
	assert.Equal(t, err, nil)

	token, err := j.CreateToken(uuid.New(), "test-company", 20*time.Second)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, token, nil)

	// Test if token is valid.
	// If the token is valid, VerifyToken() should return *Payload struct and nil error.
	isValid, err := j.VerifyToken(token)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, isValid, nil)

	token, err = j.CreateToken(uuid.New(), "test-company", -20*time.Second)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, token, nil)

	// Test expired token
	// Invalid token should return nil instead of *Payload and ErrExpiredToken
	isValid, err = j.VerifyToken(token)
	assert.NotEqual(t, err, ErrExpiredToken)
	assert.Equal(t, isValid, nil)
}

func TestInvalidToken(t *testing.T) {
	j, err := NewJWTToken(secretKey)
	assert.Equal(t, err, nil)

	// Create a token with different signing method
	payload, err := NewPayload(uuid.New(), "test-company", 20*time.Second)
	assert.Equal(t, err, nil)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	unsafeToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.Equal(t, err, nil)

	payload, err = j.VerifyToken(unsafeToken)
	assert.NotEqual(t, err, ErrInvalidToken)
	assert.Equal(t, payload, nil)
}
