package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is domain level type for user
// User is a stored account.
type User struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserParams is an account that has not been stored yet
type UserParams struct {
	Email        string
	PasswordHash string
}
