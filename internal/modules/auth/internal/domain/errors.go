package domain

import "github.com/sohan-14/auth-system-go/internal/apperr"

var (
	// ErrEmailTaken error for the already existed email registration
	ErrEmailTaken = apperr.Conflict("email_taken", "email is already registered")
	// ErrUserNotFound is returned when no user matches the lookup.
	ErrUserNotFound = apperr.NotFound("user_not_found", "User not found.")
	// ErrInvalidCredentials return the unauthorized error
	ErrInvalidCredentials = apperr.Unauthorized("invalid_credentials", "Email or password is incorrect.")
)
