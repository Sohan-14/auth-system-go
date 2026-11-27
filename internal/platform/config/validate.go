package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const maxPoolConns = 1000

// validate runs business-rule validation across all sections.
func (c *Config) validate() error {
	return errors.Join(
		c.App.validate(),
		c.HTTP.validate(),
		c.Logger.validate(c.App.Environment),
		c.Database.validate(c.App.Environment),
	)
}

func (a *App) validate() error {
	if strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("%s: must not be empty", envAppName)
	}
	return nil
}

func (h *HTTP) validate() error {
	errs := []error{
		h.validateAddr(),
		positive(envHTTPReadTimeout, h.ReadTimeout),
		positive(envHTTPReadHeaderTimeout, h.ReadHeaderTimeout),
		positive(envHTTPWriteTimeout, h.WriteTimeout),
		positive(envHTTPIdleTimeout, h.IdleTimeout),
		positive(envHTTPShutdownTimeout, h.ShutdownTimeout),
	}

	if h.DrainDelay < 0 {
		errs = append(errs, fmt.Errorf("%s: must not be negative, got %v", envHTTPDrainDelay, h.DrainDelay))
	}
	if h.MaxBodyBytes <= 0 {
		errs = append(errs, fmt.Errorf("%s: must be greater than 0, got %d", envHTTPMaxBodyBytes, h.MaxBodyBytes))
	}
	if h.MaxHeaderBytes <= 0 {
		errs = append(errs, fmt.Errorf("%s: must be greater than 0, got %d", envHTTPMaxHeaderBytes, h.MaxHeaderBytes))
	}
	if h.ReadHeaderTimeout > h.ReadTimeout {
		errs = append(errs, fmt.Errorf("%s (%v) must not exceed %s (%v)",
			envHTTPReadHeaderTimeout, h.ReadHeaderTimeout, envHTTPReadTimeout, h.ReadTimeout))
	}

	return errors.Join(errs...)
}

func (h *HTTP) validateAddr() error {
	_, port, err := net.SplitHostPort(h.Addr)
	if err != nil {
		return fmt.Errorf("%s: want host:port such as :8080: %w", envHTTPAddr, err)
	}

	// Port 0 would bind a random port, which is never intended for a service.
	if n, err := strconv.ParseUint(port, 10, 16); err != nil || n == 0 {
		return fmt.Errorf("%s: port %q must be a number from 1 to 65535", envHTTPAddr, port)
	}
	return nil
}

func (l *Logger) validate(env Environment) error {
	if env == Production && l.Format != LogFormatJSON {
		return fmt.Errorf("%s=%s: %s must be %q", envAppEnv, env, envLogFormat, LogFormatJSON)
	}
	return nil
}

func (db *Database) validate(env Environment) error {
	errs := []error{
		db.validateURL(env),
		positive(envDBMaxConnLife, db.MaxConnLifetime),
		positive(envDBMaxConnIdle, db.MaxConnIdleTime),
		positive(envDBConnectTimeout, db.ConnectTimeout),
		positive(envDBStartupTimeout, db.StartupTimeout),
		positive(envDBHealthTimeout, db.HealthTimeout),
	}

	if db.MaxConns < 1 || db.MaxConns > maxPoolConns {
		errs = append(errs, fmt.Errorf("%s: must be between 1 and %d, got %d",
			envDBMaxConns, maxPoolConns, db.MaxConns))
	}
	if db.MinConns < 0 || db.MinConns > db.MaxConns {
		errs = append(errs, fmt.Errorf("%s: must be between 0 and %s (%d), got %d",
			envDBMinConns, envDBMaxConns, db.MaxConns, db.MinConns))
	}
	if db.ConnectTimeout > db.StartupTimeout {
		errs = append(errs, fmt.Errorf("%s (%v) must not exceed %s (%v)",
			envDBConnectTimeout, db.ConnectTimeout, envDBStartupTimeout, db.StartupTimeout))
	}

	return errors.Join(errs...)
}

func (db *Database) validateURL(env Environment) error {
	u, err := url.Parse(db.URL)
	if err != nil {
		return fmt.Errorf("%s: not a valid URL", envDatabaseURL)
	}

	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return fmt.Errorf("%s: scheme must be postgres:// or postgresql://, got %q", envDatabaseURL, u.Scheme)
	}
	if strings.Trim(u.Path, "/") == "" {
		return fmt.Errorf("%s: database name is missing", envDatabaseURL)
	}

	if env == Production {
		switch sslmode := u.Query().Get("sslmode"); sslmode {
		case "require", "verify-ca", "verify-full":
		default:
			if sslmode == "" {
				sslmode = "(not set)"
			}
			return fmt.Errorf("%s: sslmode must be require, verify-ca, or verify-full when %s=%s, got %s",
				envDatabaseURL, envAppEnv, env, sslmode)
		}
	}
	return nil
}

// positive reports an error unless d is greater than zero.
func positive(key string, d time.Duration) error {
	if d <= 0 {
		return fmt.Errorf("%s: must be greater than 0, got %v", key, d)
	}
	return nil
}
