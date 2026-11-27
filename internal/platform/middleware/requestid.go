package middleware

import (
	"log/slog"
	"net/http"

	"github.com/sohan-14/auth-system-go/internal/platform/logger"
	"github.com/sohan-14/auth-system-go/internal/platform/requestid"
)

// RequestID gives every request an ID: a valid inbound X-Request-ID
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestid.Header)
		if !requestid.Valid(id) {
			id = requestid.Generate()
		}
		w.Header().Set(requestid.Header, id)

		ctx := requestid.NewContext(r.Context(), id)
		ctx = logger.ContextWithAttrs(ctx, slog.String("request_id", id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
