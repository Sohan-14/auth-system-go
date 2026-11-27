package config

import (
	"errors"
	"fmt"
	"strings"
)

// reader reads environment variables and collects parsing errors
// so every problem can be reported at once.
type reader struct {
	getenv func(string) string
	errs   []error
}

// newReader creates the reader with getenv method
func newReader(getenv func(string) string) *reader {
	return &reader{getenv: getenv}
}

// fail add the error into reader
func (r *reader) fail(key string, cause error) {
	r.errs = append(r.errs, fmt.Errorf("%s: %w", key, cause))
}

// optional reads and parses an environment variable. It returns def when the
// variable is unset or blank, and also when it is invalid (the error is
// recorded).
func optional[T any](r *reader, key string, def T, parse func(string) (T, error)) T {
	raw := strings.TrimSpace(r.getenv(key))
	if raw == "" {
		return def
	}

	v, err := parse(raw)
	if err != nil {
		r.fail(key, err)
		return def
	}
	return v
}

// required reads and parses a mandatory environment variable, recording an
// error if it is missing or invalid.
func required[T any](r *reader, key string, parse func(string) (T, error)) T {
	var zero T

	raw := strings.TrimSpace(r.getenv(key))
	if raw == "" {
		r.fail(key, errors.New("required but not set"))
		return zero
	}

	v, err := parse(raw)
	if err != nil {
		r.fail(key, err)
		return zero
	}
	return v
}
