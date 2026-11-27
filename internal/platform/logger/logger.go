// Package logger builds the application's structured slog logger.
package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/sohan-14/auth-system-go/internal/platform/config"
)

const redactedValue = "[REDACTED]"

// sensitiveParts are lowercase substrings; any attribute key containing one
// has its value replaced. This errs on the side of hiding too much.
var sensitiveParts = []string{
	"password", "passwd", "secret", "token",
	"api_key", "api-key", "apikey", "authorization", "cookie",
}

// redact is a slog ReplaceAttr hook that hides values of sensitive keys.
// It inspects keys only, never values.
func redact(_ []string, a slog.Attr) slog.Attr {
	key := strings.ToLower(a.Key)
	for _, part := range sensitiveParts {
		if strings.Contains(key, part) {
			return slog.String(a.Key, redactedValue)
		}
	}
	return a
}

// New builds a logger from cfg that writes to w (os.Stdout if nil), redacts
// sensitive keys, and includes attributes stored with ContextWithAttrs.
func New(w io.Writer, cfg *config.Logger) *slog.Logger {
	if w == nil {
		w = os.Stdout
	}

	opts := &slog.HandlerOptions{
		AddSource:   cfg.AddSource,
		Level:       cfg.Level,
		ReplaceAttr: redact,
	}

	var h slog.Handler
	if cfg.Format == config.LogFormatText {
		h = slog.NewTextHandler(w, opts)
	} else {
		h = slog.NewJSONHandler(w, opts)
	}

	return slog.New(contextHandler{h})
}

// Setup builds a logger with New and installs it as the slog default.
func Setup(w io.Writer, cfg *config.Logger) *slog.Logger {
	l := New(w, cfg)
	slog.SetDefault(l)
	return l
}
