// Package cid contains helper functions for creating
// and passing a correlation id using the context package.
// The correlation id is a UUID string used to trace
// requests and logs across the system.
package cid

import (
	"context"

	"github.com/google/uuid"
)

// Cid context key.
type cidKeyType struct{}

var cidKey cidKeyType

// New will return a new correlation id (UUID).
func New() string {
	return uuid.New().String()
}

// NewWithContext will return a new
// correlation id and inject it in the context.
func NewWithContext(ctx context.Context) (context.Context, string) {
	cid := New()
	return WithContext(ctx, cid), cid
}

// WithContext will return a new context
// which containing the correlation id.
func WithContext(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, cidKey, cid)
}

// FromContext will return the correlation id from
// the received context if it exists.
// Otherwise the function returns an empty string.
func FromContext(ctx context.Context) string {
	cid := ctx.Value(cidKey)
	if cid == nil {
		return ""
	}
	return cid.(string)
}

// FromContextOrNew will return the correlation id
// form the context if it exists.
// Otherwise it will create a new correlation id,
// inject it into the context and return it.
func FromContextOrNew(ctx context.Context) (context.Context, string) {
	if cid := FromContext(ctx); cid != "" {
		return ctx, cid
	}

	cid := New()
	return WithContext(ctx, cid), cid
}
