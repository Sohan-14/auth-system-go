package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sohan-14/auth-system-go/internal/platform/config"
)

// Server wraps http.Server with readiness draining and graceful shutdown.
type Server struct {
	srv             *http.Server
	log             *slog.Logger
	drainDelay      time.Duration
	shutdownTimeout time.Duration
	draining        atomic.Bool
	onDrain         []func()
}

// New builds a Server from cfg
func New(cfg *config.HTTP, log *slog.Logger, handler http.Handler) *Server {
	if log == nil {
		log = slog.Default()
	}
	if handler == nil {
		handler = http.NewServeMux()
	}

	return &Server{
		log:             log,
		drainDelay:      cfg.DrainDelay,
		shutdownTimeout: cfg.ShutdownTimeout,
		srv: &http.Server{
			Addr:              cfg.Addr,
			Handler:           http.MaxBytesHandler(handler, cfg.MaxBodyBytes),
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			MaxHeaderBytes:    cfg.MaxHeaderBytes,
			ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelError),
		},
	}
}

// OnDrain registers a callback run when shutdown begins, before the drain
// delay. Not safe for concurrent use: call it before Run.
func (s *Server) OnDrain(fn func()) {
	if fn != nil {
		s.onDrain = append(s.onDrain, fn)
	}
}

// IsDraining reports whether shutdown has started. Use it in /readyz.
func (s *Server) IsDraining() bool {
	return s.draining.Load()
}

// Run listens on the configured address and serves until ctx is cancelled
// (graceful shutdown, returns nil) or the server fails.
func (s *Server) Run(ctx context.Context) error {
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", s.srv.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.srv.Addr, err)
	}
	return s.serve(ctx, ln)
}

// serve is split from Run so tests can pass their own listener.
func (s *Server) serve(ctx context.Context, ln net.Listener) error {
	errCh := make(chan error, 1)
	go func() { errCh <- s.srv.Serve(ln) }()

	s.log.InfoContext(ctx, "http server listening", "addr", ln.Addr().String())

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve http: %w", err)
	case <-ctx.Done():
		// ctx is already cancelled; shutdown needs a live one.
		return s.shutdown(context.WithoutCancel(ctx))
	}
}

// shutdown flips readiness, waits for load balancers to notice, then stops
// accepting work and waits for in-flight requests
func (s *Server) shutdown(ctx context.Context) error {
	s.draining.Store(true)
	s.srv.SetKeepAlivesEnabled(false)
	for _, fn := range s.onDrain {
		fn()
	}

	s.log.InfoContext(ctx, "shutting down http server",
		"drain_delay", s.drainDelay.String(),
		"shutdown_timeout", s.shutdownTimeout.String(),
	)

	time.Sleep(s.drainDelay)

	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		_ = s.srv.Close() // force-close connections still open after the timeout
		return fmt.Errorf("graceful shutdown: %w", err)
	}

	s.log.InfoContext(ctx, "http server stopped")
	return nil
}
