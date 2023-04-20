package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const Issuer = "kperanovic"

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload implements jwt.Claims interface{} necessary for jwt claims validation.
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UserID    uuid.UUID `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload will create a new token payload
func NewPayload(id uuid.UUID, name string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Name:      name,
		UserID:    id,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p *Payload) GetIssuer() (string, error) {
	return Issuer, nil
}

func (p *Payload) GetSubject() (string, error) {
	return p.Name, nil
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return []string{p.UserID.String()}, nil
}
