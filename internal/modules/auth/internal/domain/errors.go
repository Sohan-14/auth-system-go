package domain

import "github.com/sohan-14/auth-system-go/internal/apperr"

var (
	// ErrEmailTaken error for the already existed email registration
	ErrEmailTaken = apperr.Conflict("email_taken", "email is already registered")
)
