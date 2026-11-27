package requestid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

// Header is the HTTP header used to receive and return the request ID.
const Header = "X-Request-ID"

const maxLen = 64

type ctxKey struct{}

// NewContext returns a copy of ctx carrying the request ID.
func NewContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// FromContext returns the request ID stored in ctx, or "" if none is set.
func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}

// Generate returns a new random 128-bit request ID as 32 hex characters.
func Generate() string {
	var b [16]byte
	_, _ = rand.Read(b[:]) // never fails on Go 1.24+
	return hex.EncodeToString(b[:])
}

// Valid reports whether id is safe to accept from a client: non-empty, at most
// 64 characters, and limited to letters, digits, '.', '_' and '-'. This blocks
// header and log injection.
func Valid(id string) bool {
	if id == "" || len(id) > maxLen {
		return false
	}
	for i := 0; i < len(id); i++ {
		if !isIDChar(id[i]) {
			return false
		}
	}
	return true
}

func isIDChar(c byte) bool {
	switch {
	case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
		return true
	case c == '.', c == '_', c == '-':
		return true
	}
	return false
}
