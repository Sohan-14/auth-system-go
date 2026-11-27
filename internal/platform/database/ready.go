package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

const maxReadyBackoff = 2 * time.Second

// waitReady pings until success or ctx ends, doubling the delay from
// initialDelay up to maxReadyBackoff. On failure it returns the last ping
// error, so the caller sees why the database was unreachable.
func waitReady(ctx context.Context, ping func(context.Context) error, log *slog.Logger, initialDelay time.Duration) error {
	unreachable := func(attempt int, err error) error {
		return fmt.Errorf("database not reachable after %d attempts: %w", attempt, err)
	}

	delay := initialDelay
	for attempt := 1; ; attempt++ {
		err := ping(ctx)
		if err == nil {
			return nil
		}
		if ctx.Err() != nil {
			return unreachable(attempt, err)
		}

		log.WarnContext(ctx, "database not ready",
			"attempt", attempt,
			"retry_in", delay.String(),
			"error", err,
		)

		select {
		case <-ctx.Done():
			return unreachable(attempt, err)
		case <-time.After(delay):
		}
		delay = min(delay*2, maxReadyBackoff)
	}
}
