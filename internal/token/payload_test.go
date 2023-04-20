package token

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

func TestNewPayload(t *testing.T) {
	payload, err := NewPayload(uuid.New(), "test", 20*time.Second)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, payload, nil)
}
