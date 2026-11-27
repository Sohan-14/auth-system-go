package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/sohan-14/auth-system-go/internal/apperr"
)

// Recover turns a handler panic into a logged 500 response
func Recover(log *slog.Logger, onError ErrorHandler) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &statusWriter{ResponseWriter: w}

			defer func() {
				rec := recover()
				if rec == nil {
					return
				}
				if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					panic(rec) // deliberate abort; let net/http handle it
				}

				log.ErrorContext(r.Context(), "panic recovered",
					"panic", fmt.Sprint(rec),
					"stack", string(debug.Stack()),
				)
				if sw.wrote() {
					return
				}
				onError(sw, r, apperr.Internal("panic", fmt.Errorf("%v", rec)))
			}()

			next.ServeHTTP(sw, r)
		})
	}
}
