package apperr

import (
	"errors"
	"fmt"
)

// FieldError describes a single input field failure during validation.
type FieldError struct {
	Field   string
	Code    string
	Message string
}

// Error is the unified structured error returned across domain, usecase, and
// persistence boundaries.
//
// Message is safe to expose to external clients for 4xx-equivalent kinds.
// cause holds the underlying internal cause for server-side logging and errors.Is/As.
type Error struct {
	Kind    Kind
	Code    string
	Message string
	Fields  []FieldError
	cause   error
}

// Error implements the standard error interface for internal logging.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	switch {
	case e.cause != nil && e.Code != "":
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.cause)
	case e.cause != nil:
		return fmt.Sprintf("%s: %v", e.Message, e.cause)
	case e.Code != "":
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	default:
		return e.Message
	}
}

// Unwrap exposes the underlying error so errors.Is and errors.As keep working
// through wrapped domain and repository errors.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// Is lets errors.Is match two *Error values by Kind and non-empty Code.
// Only the target itself is inspected; chain traversal is done by errors.Is.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok || e == nil || t == nil {
		return false
	}
	return t.Code != "" && e.Kind == t.Kind && e.Code == t.Code
}

// New creates a structured Error with a semantic Kind,
// machine-readable Code, and client-safe Message.
func New(kind Kind, code, message string) *Error {
	return &Error{
		Kind:    kind,
		Code:    code,
		Message: message,
	}
}

// Wrap attaches an underlying technical error (such as a pgx or stdlib error)
// for internal logging while preserving a safe client-facing message.
func Wrap(kind Kind, code, message string, cause error) *Error {
	return &Error{
		Kind:    kind,
		Code:    code,
		Message: message,
		cause:   cause,
	}
}

// As returns the first *Error in err's chain, if any.
func As(err error) (*Error, bool) {
	return errors.AsType[*Error](err)
}

// KindOf extracts the semantic Kind from any error. It returns "" for nil and
// KindInternal for plain or unrecognized errors.
func KindOf(err error) Kind {
	if err == nil {
		return ""
	}
	if appErr, ok := As(err); ok {
		return appErr.Kind
	}
	return KindInternal
}
