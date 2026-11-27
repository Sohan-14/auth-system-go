package httpx

const (
	defaultProblemType  = "about:blank"
	genericServerDetail = "An unexpected internal error occurred."
)

// FieldError describes a single invalid request field returned on validation
// failures so clients can highlight specific inputs.
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ProblemDetail represents an RFC 9457 Problem Details payload serialized as
// application/problem+json. Standard fields (type, title, status, detail,
// instance) pair with project extensions (code, request_id, errors) for
// client logic and log correlation.
type ProblemDetail struct {
	Type      string       `json:"type"`
	Title     string       `json:"title"`
	Status    int          `json:"status"`
	Detail    string       `json:"detail,omitempty"`
	Instance  string       `json:"instance,omitempty"`
	Code      string       `json:"code,omitempty"`
	RequestID string       `json:"request_id,omitempty"`
	Errors    []FieldError `json:"errors,omitempty"`
}
