package middleware

import (
	"net/http"

	"github.com/sohan-14/auth-system-go/internal/apperr"
)

// MaxBodySize limits request bodies to limit bytes. A larger declared
// Content-Length is rejected immediately with 413; a body that turns out larger
// while being read (chunked or lying client) fails the read with
// *http.MaxBytesError, which httpx.DecodeJSON maps to 413 as well.
func MaxBodySize(limit int64, onError ErrorHandler) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > limit {
				onError(w, r, apperr.TooLarge("request_body_too_large", "Request body exceeds maximum allowed size."))
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}
