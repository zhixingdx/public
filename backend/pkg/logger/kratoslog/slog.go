package log

import (
	"context"
	"log/slog"
)

// AsSlog bridges the Kratos v2-style logger used by existing business code to
// Kratos v3's slog-based APIs.
func AsSlog(logger Logger) *slog.Logger {
	if logger == nil {
		logger = DefaultLogger
	}
	return slog.New(&slogHandler{logger: logger})
}

type slogHandler struct {
	logger Logger
	attrs  []slog.Attr
	group  string
}

func (h *slogHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *slogHandler) Handle(ctx context.Context, record slog.Record) error {
	level := LevelInfo
	switch {
	case record.Level <= slog.LevelDebug:
		level = LevelDebug
	case record.Level >= slog.LevelError:
		level = LevelError
	case record.Level >= slog.LevelWarn:
		level = LevelWarn
	}

	keyvals := make([]any, 0, 2+len(h.attrs)*2+record.NumAttrs()*2)
	keyvals = append(keyvals, DefaultMessageKey, record.Message)
	appendAttr := func(a slog.Attr) bool {
		a.Value = a.Value.Resolve()
		key := a.Key
		if h.group != "" {
			key = h.group + "." + key
		}
		keyvals = append(keyvals, key, a.Value.Any())
		return true
	}
	for _, attr := range h.attrs {
		appendAttr(attr)
	}
	record.Attrs(appendAttr)
	return WithContext(ctx, h.logger).Log(level, keyvals...)
}

func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := *h
	next.attrs = append(append([]slog.Attr{}, h.attrs...), attrs...)
	return &next
}

func (h *slogHandler) WithGroup(name string) slog.Handler {
	next := *h
	if next.group == "" {
		next.group = name
	} else {
		next.group += "." + name
	}
	return &next
}
