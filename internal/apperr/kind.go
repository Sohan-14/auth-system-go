package apperr

// Kind classifies an application error semantically without coupling lower
// layers (domain, usecase, store) to HTTP status codes or transport details.
type Kind string

const (
	// KindInvalid is for malformed syntax or unreadable request payloads.
	KindInvalid Kind = "invalid"
	// KindValidation is for syntactically valid input that violates domain rules.
	KindValidation Kind = "validation"
	// KindUnauthorized is for missing, expired, or invalid authentication.
	KindUnauthorized Kind = "unauthorized"
	// KindForbidden is for an authenticated caller lacking permission.
	KindForbidden Kind = "forbidden"
	// KindNotFound is for a requested resource that does not exist.
	KindNotFound Kind = "not_found"
	// KindConflict is for a state conflict, such as a duplicate unique key.
	KindConflict Kind = "conflict"
	// KindTooLarge is for a request with more data than expected.
	KindTooLarge Kind = "too_large"
	// KindUnsupportedMedia is for an unexpected media type (e.g. text sent to a JSON endpoint).
	KindUnsupportedMedia Kind = "unsupported_media"
	// KindRateLimited is for callers exceeding allowed request thresholds.
	KindRateLimited Kind = "rate_limited"
	// KindUnavailable is for a service that cannot serve requests right now.
	KindUnavailable Kind = "unavailable"
	// KindInternal is for unexpected server, database, or runtime failure.
	KindInternal Kind = "internal"
)
