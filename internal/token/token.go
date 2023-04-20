package token

import (
	"time"

	"github.com/google/uuid"
)

type Token interface {
	// CreateToken creates a new token for a specific name and duration.
	CreateToken(id uuid.UUID, name string, duration time.Duration) (string, error)
	// Verify token checks if the provided token is valid.
	VerifyToken(token string) (*Payload, error)
}
