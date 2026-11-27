package apperr

// Internal wraps an unexpected technical failure as a KindInternal error.
// The wrapped err is logged internally and never sent to the client.
func Internal(message string, err error) *Error {
	return Wrap(KindInternal, "internal_error", message, err)
}

// Invalid constructs a KindInvalid error for malformed input or bad syntax.
func Invalid(code, message string) *Error {
	return New(KindInvalid, code, message)
}

// Validation constructs a KindValidation error, optionally carrying field-level
// details for form or payload validation responses.
func Validation(message string, fields ...FieldError) *Error {
	return &Error{
		Kind:    KindValidation,
		Code:    "validation_error",
		Message: message,
		Fields:  fields,
	}
}

// Unauthorized constructs a KindUnauthorized error.
func Unauthorized(code, message string) *Error {
	return New(KindUnauthorized, code, message)
}

// Forbidden constructs a KindForbidden error.
func Forbidden(code, message string) *Error {
	return New(KindForbidden, code, message)
}

// NotFound constructs a KindNotFound error.
func NotFound(code, message string) *Error {
	return New(KindNotFound, code, message)
}

// Conflict constructs a KindConflict error.
func Conflict(code, message string) *Error {
	return New(KindConflict, code, message)
}

// TooLarge constructs a KindTooLarge error.
func TooLarge(code, message string) *Error {
	return New(KindTooLarge, code, message)
}

// UnsupportedMedia constructs a KindUnsupportedMedia error.
func UnsupportedMedia(code, message string) *Error {
	return New(KindUnsupportedMedia, code, message)
}

// RateLimited constructs a KindRateLimited error.
func RateLimited(code, message string) *Error {
	return New(KindRateLimited, code, message)
}

// Unavailable constructs a KindUnavailable error.
func Unavailable(code, message string) *Error {
	return New(KindUnavailable, code, message)
}
