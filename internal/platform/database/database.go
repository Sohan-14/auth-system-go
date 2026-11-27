package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sohan-14/auth-system-go/internal/platform/config"
)

const readyInitialDelay = 250 * time.Millisecond

// Database wraps a pgx connection pool with startup and health checks.
type Database struct {
	pool          *pgxpool.Pool
	log           *slog.Logger
	healthTimeout time.Duration
}

// New builds the pool and blocks until the database answers a ping, or
// cfg.StartupTimeout (or ctx) expires
func New(ctx context.Context, cfg *config.Database, log *slog.Logger) (*Database, error) {
	if log == nil {
		log = slog.Default()
	}

	poolCfg, err := buildPoolConfig(cfg)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	readyCtx, cancel := context.WithTimeout(ctx, cfg.StartupTimeout)
	defer cancel()

	if err := waitReady(readyCtx, pool.Ping, log, readyInitialDelay); err != nil {
		pool.Close()
		return nil, err
	}

	log.InfoContext(ctx, "database connected",
		"host", poolCfg.ConnConfig.Host,
		"port", poolCfg.ConnConfig.Port,
		"database", poolCfg.ConnConfig.Database,
		"max_conns", poolCfg.MaxConns,
		"min_conns", poolCfg.MinConns,
	)

	return &Database{pool: pool, log: log, healthTimeout: cfg.HealthTimeout}, nil
}

// Pool returns the underlying pool; it satisfies sqlc generated DBTX.
func (db *Database) Pool() *pgxpool.Pool {
	return db.pool
}

// Health pings the database within the configured health timeout.
// Use it for readiness (/readyz) only, never for liveness.
func (db *Database) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, db.healthTimeout)
	defer cancel()

	if err := db.pool.Ping(ctx); err != nil {
		return fmt.Errorf("database health check: %w", err)
	}
	return nil
}

// Close releases all connections. Call it after in-flight handlers finish.
func (db *Database) Close() {
	db.log.Info("closing database pool")
	db.pool.Close()
}

// buildPoolConfig parses the URL and applies pool settings from cfg.
func buildPoolConfig(cfg *config.Database) (*pgxpool.Config, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnLifetimeJitter = cfg.MaxConnLifetime / 10
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

	return poolCfg, nil
}
