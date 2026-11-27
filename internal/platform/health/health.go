// Package health provides liveness and readiness HTTP handlers.
package health

import (
	"context"
	"log/slog"
	"net/http"
	"sync/atomic"
)

// CheckFunc reports whether a dependency is usable; nil means healthy.
type CheckFunc func(ctx context.Context) error

type check struct {
	name string
	fn   CheckFunc
}

// Handler serves liveness and readiness probes.
type Handler struct {
	log      *slog.Logger
	checks   []check
	draining atomic.Bool
}

// New creates a Handler. A nil log uses slog.Default().
func New(log *slog.Logger) *Handler {
	if log == nil {
		log = slog.Default()
	}
	return &Handler{log: log}
}

// Add registers a readiness check. Call it before the server starts.
// Each check must bound its own runtime (db.Health already does).
func (h *Handler) Add(name string, fn CheckFunc) {
	h.checks = append(h.checks, check{name, fn})
}

// Drain makes readiness fail from now on. Register it with Server.OnDrain.
func (h *Handler) Drain() {
	h.draining.Store(true)
}

// Live reports that the process is running. It checks no dependencies, so a
// database outage never gets a healthy pod restarted.
func (h *Handler) Live(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// Ready returns 204 when the service can take traffic, and 503 while
// draining or when any check fails. Failure details are logged, not returned.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if h.draining.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	for _, c := range h.checks {
		if err := c.fn(r.Context()); err != nil {
			h.log.WarnContext(r.Context(), "readiness check failed", "check", c.name, "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
