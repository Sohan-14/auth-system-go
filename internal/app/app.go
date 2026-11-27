package app

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/sohan-14/auth-system-go/internal/platform/config"
	"github.com/sohan-14/auth-system-go/internal/platform/database"
	"github.com/sohan-14/auth-system-go/internal/platform/health"
	"github.com/sohan-14/auth-system-go/internal/platform/httpserver"
	"github.com/sohan-14/auth-system-go/internal/platform/httpx"
	"github.com/sohan-14/auth-system-go/internal/platform/logger"
	"github.com/sohan-14/auth-system-go/internal/platform/middleware"
	"github.com/sohan-14/auth-system-go/internal/platform/router"
)

// Run builds every dependency, serves until SIGINT/SIGTERM, then shuts down in
// reverse order. Deferred calls run before main sees the returned error.
func Run(ctx context.Context,
	getenv func(string) string,
	stdout io.Writer) error {
	// config
	cfg, err := config.Load(getenv)
	if err != nil {
		return err
	}

	// logger
	log := logger.Setup(stdout, &cfg.Logger)
	log.Info("starting", "app", cfg.App.Name, "env", cfg.App.Environment)

	// graceful shutdown
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() { // after the first signal, a second Ctrl+C kills the process
		<-ctx.Done()
		stop()
	}()

	// database
	db, err := database.New(ctx, &cfg.Database, log)
	if err != nil {
		return err
	}
	defer db.Close()

	hc := health.New(log)
	hc.Add("database", db.Health)

	// router
	r := router.New()

	// Operational endpoints: unversioned and unauthenticated.
	r.HandleFunc("GET /healthz", hc.Live)
	r.HandleFunc("GET /readyz", hc.Ready)

	// global middleware
	handler := middleware.Chain(r,
		middleware.RequestID,
		middleware.AccessLog(log),
		middleware.Recover(log, httpx.NewErrorWriter(log).Write),
		middleware.MaxBodySize(cfg.HTTP.MaxBodyBytes, httpx.NewErrorWriter(log).Write),
	)

	// server
	srv := httpserver.New(&cfg.HTTP, log, handler)
	srv.OnDrain(hc.Drain)
	if err := srv.Run(ctx); err != nil {
		return err
	}

	log.Info("stopped")
	return nil
}
