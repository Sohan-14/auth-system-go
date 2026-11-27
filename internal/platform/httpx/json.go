package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/sohan-14/auth-system-go/internal/apperr"
)

// DecodeJSON decodes a single top-level JSON value from r.Body into dst.
// Unknown fields and trailing data are rejected. Wrap r.Body with
// http.MaxBytesReader (middleware) before calling.
func DecodeJSON(r *http.Request, dst any) error {
	if r == nil || r.Body == nil || r.Body == http.NoBody {
		return apperr.Invalid("empty_body", "Request body must not be empty.")
	}

	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mt != "application/json" {
		return apperr.UnsupportedMedia("unsupported_media_type", "Content-Type must be application/json.")
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return classifyDecodeError(err)
	}

	// Reject trailing tokens or multiple JSON values.
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			return tooLarge()
		}
		return apperr.Invalid("trailing_json", "Request body must contain only a single JSON value.")
	}
	return nil
}

func tooLarge() *apperr.Error {
	return apperr.TooLarge("request_body_too_large", "Request body exceeds maximum allowed size.")
}

func classifyDecodeError(err error) error {
	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return tooLarge()
	}
	if errors.Is(err, io.EOF) {
		return apperr.Invalid("empty_body", "Request body must not be empty.")
	}
	if _, ok := errors.AsType[*json.SyntaxError](err); ok || errors.Is(err, io.ErrUnexpectedEOF) {
		return apperr.Invalid("malformed_json", "Request body contains invalid JSON.")
	}
	if e, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
		if e.Field != "" {
			return apperr.Invalid("invalid_field_type", fmt.Sprintf("Invalid value for field %q.", e.Field))
		}
		return apperr.Invalid("invalid_field_type", "Invalid value type in request body.")
	}
	// DisallowUnknownFields produces a plain error: `json: unknown field "x"`.
	if strings.HasPrefix(err.Error(), "json: unknown field ") {
		return apperr.Invalid("unknown_field", "Request body contains an unknown field.")
	}
	return apperr.Invalid("invalid_request", "Unable to parse request payload.")
}

// WriteJSON encodes payload as JSON and writes it with the given status.
// Encoding happens first, so a failure leaves the response untouched and the
// caller can delegate to ErrorWriter.
func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return fmt.Errorf("encode json response: %w", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
	return nil
}
