package domain

import (
	"net/mail"
	"strings"

	"github.com/sohan-14/auth-system-go/internal/apperr"
)

// Email is a normalized (trimmed, lowercase) email address.
type Email string

const maxEmailLen = 254 // RFC 5321 path limit

func parseEmail(email string) (Email, *apperr.FieldError) {
	email = strings.ToLower(strings.TrimSpace(email))

	switch {
	case email == "":
		return "", &apperr.FieldError{Field: "email", Code: "required", Message: "Email is required."}
	case len(email) > maxEmailLen:
		return "", &apperr.FieldError{Field: "email", Code: "too_long", Message: "Email is too long."}
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email || !strings.Contains(email[strings.LastIndex(email, "@"):], ".") {
		return "", &apperr.FieldError{Field: "email", Code: "invalid", Message: "Email must be a valid address."}
	}
	return Email(email), nil
}
