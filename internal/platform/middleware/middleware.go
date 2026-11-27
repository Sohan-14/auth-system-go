package middleware

import (
	"net/http"
	"slices"
)

// Middleware wraps a handler with extra behavior.
type Middleware func(http.Handler) http.Handler

// ErrorHandler renders err as an HTTP response (e.g. httpx.ErrorWriter.Write).
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

// Chain wraps h(handler) with mws.
// The first middleware is the outermost and runs first.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for _, mw := range slices.Backward(mws) {
		if mw != nil {
			h = mw(h)
		}
	}
	return h
}

// statusWriter records the status code and body size of a response.
type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int64
}

// WriteHeader records the first non-informational
// HTTP status code and delegates to the underlying writer.
func (w *statusWriter) WriteHeader(code int) {
	if w.status == 0 && code >= http.StatusOK { // ignore 1xx informational
		w.status = code
	}
	w.ResponseWriter.WriteHeader(code)
}

// Write records a default 200 OK status
// if needed, writes the payload, and tracks byte count.
func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += int64(n)
	return n, err
}

// Unwrap lets http.ResponseController reach Flush, Hijack, etc.
func (w *statusWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *statusWriter) wrote() bool {
	return w.status != 0
}

func (w *statusWriter) statusOrOK() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}
