package config

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// Load reads and validates configuration in two phases:
//  1. Read every environment variable, aggregating all parse errors.
//  2. Validate the result against business rules, aggregating all violations.
func Load(getenv func(string) string) (Config, error) {
	r := newReader(getenv)

	cfg := Config{
		App:      loadApp(r),
		HTTP:     loadHTTP(r),
		Logger:   loadLogger(r),
		Database: loadDatabase(r),
	}

	if err := errors.Join(r.errs...); err != nil {
		return Config{}, fmt.Errorf("config:\n%w", err)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("config:\n%w", err)
	}

	return cfg, nil
}

func loadApp(r *reader) App {
	return App{
		Name:        optional(r, envAppName, "auth-system-go", asString),
		Environment: required(r, envAppEnv, oneOf(Development, Staging, Production)),
	}
}

func loadHTTP(r *reader) HTTP {
	return HTTP{
		Addr:              required(r, envHTTPAddr, asString),
		ReadTimeout:       optional(r, envHTTPReadTimeout, 15*time.Second, asDuration),
		ReadHeaderTimeout: optional(r, envHTTPReadHeaderTimeout, 5*time.Second, asDuration),
		WriteTimeout:      optional(r, envHTTPWriteTimeout, 30*time.Second, asDuration),
		IdleTimeout:       optional(r, envHTTPIdleTimeout, 60*time.Second, asDuration),
		ShutdownTimeout:   optional(r, envHTTPShutdownTimeout, 20*time.Second, asDuration),
		DrainDelay:        optional(r, envHTTPDrainDelay, 0, asDuration),
		MaxBodyBytes:      optional(r, envHTTPMaxBodyBytes, 1<<20, asInt64),
		MaxHeaderBytes:    optional(r, envHTTPMaxHeaderBytes, 64<<10, asInt),
	}
}

func loadLogger(r *reader) Logger {
	return Logger{
		Level:     optional(r, envLogLevel, slog.LevelInfo, asLevel),
		Format:    optional(r, envLogFormat, LogFormatJSON, oneOf(LogFormatJSON, LogFormatText)),
		AddSource: optional(r, envLogAddSource, false, asBool),
	}
}

func loadDatabase(r *reader) Database {
	return Database{
		URL:             required(r, envDatabaseURL, asString),
		MaxConns:        optional(r, envDBMaxConns, 10, asInt32),
		MinConns:        optional(r, envDBMinConns, 2, asInt32),
		MaxConnLifetime: optional(r, envDBMaxConnLife, 30*time.Minute, asDuration),
		MaxConnIdleTime: optional(r, envDBMaxConnIdle, 5*time.Minute, asDuration),
		ConnectTimeout:  optional(r, envDBConnectTimeout, 5*time.Second, asDuration),
		StartupTimeout:  optional(r, envDBStartupTimeout, 15*time.Second, asDuration),
		HealthTimeout:   optional(r, envDBHealthTimeout, 3*time.Second, asDuration),
	}
}
