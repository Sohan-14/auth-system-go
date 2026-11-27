package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/sohan-14/auth-system-go/internal/apperr"
	"github.com/sohan-14/auth-system-go/internal/platform/requestid"
)

// ErrorWriter translates errors into RFC 9457 responses and logs server-side
// failures with request correlation.
type ErrorWriter struct {
	logger *slog.Logger
}

// NewErrorWriter constructs an ErrorWriter.
func NewErrorWriter(logger *slog.Logger) *ErrorWriter {
	if logger == nil {
		logger = slog.Default()
	}
	return &ErrorWriter{logger: logger}
}

// Write renders err as an application/problem+json response
func (ew *ErrorWriter) Write(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil || errors.Is(err, context.Canceled) {
		return
	}

	reqId := requestid.FromContext(r.Context())

	problem := buildProblem(err, reqId, r.URL.Path)

	if problem.Status >= http.StatusInternalServerError {
		ew.logger.ErrorContext(r.Context(), "request failed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("request_id", reqId),
			slog.Int("status", problem.Status),
			slog.Any("error", err),
		)
	}

	h := w.Header()
	h.Set("Content-Type", "application/problem+json")
	h.Set("Cache-Control", "no-store")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(problem.Status)

	if encErr := json.NewEncoder(w).Encode(problem); encErr != nil {
		ew.logger.ErrorContext(r.Context(), "failed to encode problem detail",
			slog.String("request_id", reqId),
			slog.Any("error", encErr),
		)
	}
}

// buildProblem maps any error to the response a client is allowed to see.
// It defaults to an opaque 500 so nothing internal leaks by accident.
func buildProblem(err error, reqId string, instance string) ProblemDetail {
	p := ProblemDetail{
		Type:      defaultProblemType,
		Status:    http.StatusInternalServerError,
		Code:      "internal_error",
		Detail:    genericServerDetail,
		RequestID: reqId,
		Instance:  instance,
	}

	if appErr, ok := errors.AsType[*apperr.Error](err); ok && appErr.Kind != apperr.KindInternal {
		p.Status = mapKindToStatus(appErr.Kind)
		p.Code = appErr.Code

		// Messages are client-safe for 4xx and Unavailable; other 5xx stay generic.
		if p.Status < http.StatusInternalServerError || appErr.Kind == apperr.KindUnavailable {
			p.Detail = appErr.Message
			p.Errors = mapFieldErrors(appErr.Fields)
		}
	} else if errors.Is(err, context.DeadlineExceeded) {
		p.Status = http.StatusGatewayTimeout
		p.Code = "request_timeout"
		p.Detail = "The request timed out."
	}

	p.Title = http.StatusText(p.Status)
	return p
}

func mapFieldErrors(fields []apperr.FieldError) []FieldError {
	if len(fields) == 0 {
		return nil
	}
	out := make([]FieldError, len(fields))
	for i, f := range fields {
		out[i] = FieldError{Field: f.Field, Code: f.Code, Message: f.Message}
	}
	return out
}
