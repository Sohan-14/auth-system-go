package domain

import "github.com/sohan-14/auth-system-go/internal/apperr"

// ParseRegistration validates registration input and returns the normalized
// email. All field problems are reported together.
func ParseRegistration(rawEmail, password string) (Email, error) {
	email, emailErr := parseEmail(rawEmail)
	passErr := checkPassword(password)

	var fields []apperr.FieldError
	if emailErr != nil {
		fields = append(fields, *emailErr)
	}
	if passErr != nil {
		fields = append(fields, *passErr)
	}
	if len(fields) > 0 {
		return "", apperr.Validation("Invalid registration data.", fields...)
	}
	return email, nil
}
