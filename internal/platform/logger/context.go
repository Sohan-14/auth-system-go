package logger

import (
	"context"
	"log/slog"
	"slices"
)

type ctxKey struct{}

// ContextWithAttrs returns a context carrying attrs, which every *Context log
// call using it will include. An attr replaces an existing one with the same
// key, so keys never repeat.
func ContextWithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	if len(attrs) == 0 {
		return ctx
	}

	existing := attrsFrom(ctx)
	merged := make([]slog.Attr, 0, len(existing)+len(attrs))
	for _, e := range existing {
		replaced := slices.ContainsFunc(attrs, func(a slog.Attr) bool { return a.Key == e.Key })
		if !replaced {
			merged = append(merged, e)
		}
	}
	merged = append(merged, attrs...)

	return context.WithValue(ctx, ctxKey{}, merged)
}

func attrsFrom(ctx context.Context) []slog.Attr {
	attrs, _ := ctx.Value(ctxKey{}).([]slog.Attr)
	return attrs
}

// contextHandler adds the attributes stored in the context to each record.
type contextHandler struct {
	slog.Handler
}

// Handle implements slog.Handler.
//
//nolint:gocritic // slog.Handler requires slog.Record by value.
func (h contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs := attrsFrom(ctx); len(attrs) > 0 {
		r = r.Clone() // required before modifying a record
		r.AddAttrs(attrs...)
	}
	return h.Handler.Handle(ctx, r)
}

// WithAttrs implements slog.Handler.
func (h contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return contextHandler{h.Handler.WithAttrs(attrs)}
}

// WithGroup implements slog.Handler.
func (h contextHandler) WithGroup(name string) slog.Handler {
	return contextHandler{h.Handler.WithGroup(name)}
}
