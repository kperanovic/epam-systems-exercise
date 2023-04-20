package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type JWTToken struct {
	secretKey string
}

// NewJWTToken will create a create a new JWTToken{} struct with secretKey.
func NewJWTToken(secretKey string) (*JWTToken, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: key must have at least %d characters", minSecretKeySize)
	}

	return &JWTToken{
		secretKey: secretKey,
	}, nil
}

// CreateToken will create a new Payload{} struct with the given inputs.
// Token is then signed with jwt.SigningMethodHS256.
// Returns a complete,signed JWT.
func (j *JWTToken) CreateToken(id uuid.UUID, name string, duration time.Duration) (string, error) {
	payload, err := NewPayload(id, name, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(j.secretKey))
}

// Verify token checks if the provided token is valid.
func (j *JWTToken) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
