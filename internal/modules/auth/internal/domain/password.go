package domain

import (
	"fmt"
	"unicode/utf8"

	"github.com/sohan-14/auth-system-go/internal/apperr"
)

const (
	minPasswordLen = 8
	maxPasswordLen = 128
)

// checkPassword is never trimmed: spaces are legitimate password characters.
func checkPassword(password string) *apperr.FieldError {
	n := utf8.RuneCountInString(password)

	switch {
	case n == 0:
		return &apperr.FieldError{Field: "password", Code: "required", Message: "Password is required."}
	case n < minPasswordLen:
		return &apperr.FieldError{Field: "password", Code: "too_short",
			Message: fmt.Sprintf("Password must be at least %d characters.", minPasswordLen)}
	case n > maxPasswordLen:
		return &apperr.FieldError{Field: "password", Code: "too_long",
			Message: fmt.Sprintf("Password must be at most %d characters.", maxPasswordLen)}
	}
	return nil
}
