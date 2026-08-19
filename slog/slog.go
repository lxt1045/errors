package slog

import (
	"context"
	"log/slog"
	"time"

	"github.com/lxt1045/errors"
	"github.com/lxt1045/errors/zerolog"
	zlog "github.com/rs/zerolog"
)

// handler 是一个自定义的 handler 包装器，用于跳过调用栈
type handler struct {
	slog.Handler
}

func toHandler(h slog.Handler) slog.Handler {
	if _, ok := h.(*handler); h == nil || ok {
		return h
	}

	return &handler{
		Handler: h,
	}
}

func NewHandler(h slog.Handler) slog.Handler {
	if h == nil {
		return nil
	}
	return &handler{
		Handler: h,
	}
}

func New(h slog.Handler) *slog.Logger {
	if h == nil {
		return slog.New(h)
	}
	return slog.New(&handler{
		Handler: h,
	})
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	if r.PC == 0 {
		f := errors.CallersSkip(4 - 2)[0]
		src := &slog.Source{
			Function: f.Func,
			File:     f.File,
			Line:     f.Line,
		}
		r.AddAttrs(slog.Attr{slog.SourceKey, slog.AnyValue(src)})
		return h.Handler.Handle(ctx, r)
	}
	f := errors.CallerFrame(r.PC)
	src := &slog.Source{
		Function: f.Func,
		File:     f.File,
		Line:     f.Line,
	}
	r.PC = 0
	r.AddAttrs(slog.Attr{slog.SourceKey, slog.AnyValue(src)})
	return h.Handler.Handle(ctx, r)
}

func (h *handler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.Handler.Enabled(ctx, l)
}

func (h *handler) WithAttrs(as []slog.Attr) slog.Handler {
	return toHandler(h.Handler.WithAttrs(as))
}

func (h *handler) WithGroup(name string) slog.Handler {
	return toHandler(h.Handler.WithGroup(name))
}

// handler 是一个自定义的 handler 包装器，用于跳过调用栈
type loggerHandler struct {
	zerolog.Logger
	prefix string // group prefix for nested groups
	attrs  []slog.Attr
}

func NewLoggerHandler(l zerolog.Logger) *loggerHandler {
	return &loggerHandler{
		Logger: l,
	}
}
func (h *loggerHandler) Enabled(_ context.Context, level slog.Level) bool {
	zl := LevelFromSlog(level)
	if zl < zerolog.GlobalLevel() {
		return false
	}
	return zl >= h.Logger.GetLevel()
}

func (h *loggerHandler) Handle(ctx context.Context, r slog.Record) error {
	zlevel := LevelFromSlog(r.Level)
	event := h.Logger.WithLevel(zlevel)
	if event == nil {
		return nil
	}

	// Propagate slog context to the zerolog event so that hooks
	// relying on Event.GetCtx() (e.g. tracing) can access it.
	if ctx != nil {
		event = event.Ctx(ctx)
	}

	var src *slog.Source
	if r.PC == 0 {
		f := errors.CallersSkip(4 - 2)[0]
		src = &slog.Source{
			Function: f.Func,
			File:     f.File,
			Line:     f.Line,
		}
		event.Any(slog.SourceKey, src)
	} else {
		f := errors.CallerFrame(r.PC)
		src = &slog.Source{
			Function: f.Func,
			File:     f.File,
			Line:     f.Line,
		}
		r.PC = 0
		event.Any(slog.SourceKey, src)
	}

	// Add pre-attached attrs from WithAttrs
	for _, a := range h.attrs {
		event = appendSlogAttr(event, a, h.prefix)
	}

	// Add attrs from the record itself
	r.Attrs(func(a slog.Attr) bool {
		event = appendSlogAttr(event, a, h.prefix)
		return true
	})

	// Add timestamp from the slog record, but only if the logger doesn't
	// already have a timestampHook (added via .With().Timestamp()) to
	// avoid duplicate timestamp keys in the output.
	if !r.Time.IsZero() {
		event.Time(zlog.TimestampFieldName, r.Time)
	}

	event.Msg(r.Message)
	return nil
}

// WithAttrs returns a new Handler with the given attributes pre-attached.
// These attributes will be included in every subsequent log record.
func (h *loggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := h.clone()
	h2.attrs = append(h2.attrs, attrs...)
	return h2
}

// WithGroup returns a new Handler with the given group name. All subsequent
// attributes will be nested under this group name in the output.
func (h *loggerHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	if h2.prefix != "" {
		h2.prefix = h2.prefix + "." + name
	} else {
		h2.prefix = name
	}
	return h2
}

func (h *loggerHandler) clone() *loggerHandler {
	h2 := &loggerHandler{
		Logger: h.Logger,
		prefix: h.prefix,
	}
	if len(h.attrs) > 0 {
		h2.attrs = make([]slog.Attr, len(h.attrs))
		copy(h2.attrs, h.attrs)
	}
	return h2
}

// LevelFromSlog maps slog levels to zerolog levels.
//
// slog levels:  Debug=-4, Info=0, Warn=4, Error=8
// zerolog levels: Trace=-1, Debug=0, Info=1, Warn=2, Error=3, Fatal=4, Panic=5
func LevelFromSlog(level slog.Level) zerolog.Level {
	switch {
	case level < slog.LevelDebug:
		return zerolog.TraceLevel
	case level < slog.LevelInfo:
		return zerolog.DebugLevel
	case level < slog.LevelWarn:
		return zerolog.InfoLevel
	case level < slog.LevelError:
		return zerolog.WarnLevel
	default:
		return zerolog.ErrorLevel
	}
}

// zerologToSlogLevel maps zerolog levels to slog levels.
func zerologToSlogLevel(level zerolog.Level) slog.Level {
	switch level {
	case zerolog.TraceLevel:
		return slog.LevelDebug - 4
	case zerolog.DebugLevel:
		return slog.LevelDebug
	case zerolog.InfoLevel:
		return slog.LevelInfo
	case zerolog.WarnLevel:
		return slog.LevelWarn
	case zerolog.ErrorLevel:
		return slog.LevelError
	case zerolog.FatalLevel:
		return slog.LevelError + 4
	case zerolog.PanicLevel:
		return slog.LevelError + 8
	default:
		return slog.LevelInfo
	}
}

// joinPrefix concatenates a prefix and key with a dot separator.
// It avoids allocations when either prefix or key is empty.
func joinPrefix(prefix, key string) string {
	if prefix == "" {
		return key
	}
	if key == "" {
		return prefix
	}
	return prefix + "." + key
}

// appendSlogAttr appends a single slog.Attr to the zerolog event, handling
// type-specific encoding to avoid reflection where possible.
func appendSlogAttr(event *zerolog.Event, attr slog.Attr, prefix string) *zerolog.Event {
	if event == nil {
		return event
	}

	// Resolve the attribute to handle LogValuer types.
	// This handles slog.KindLogValuer implicitly by unwrapping
	// any values that implement slog.LogValuer to their resolved form.
	attr.Value = attr.Value.Resolve()

	// For group kinds, handle grouping before key concatenation
	if attr.Value.Kind() == slog.KindGroup {
		attrs := attr.Value.Group()
		if len(attrs) == 0 {
			return event
		}
		groupPrefix := joinPrefix(prefix, attr.Key)
		for _, ga := range attrs {
			event = appendSlogAttr(event, ga, groupPrefix)
		}
		return event
	}

	// Skip empty keys for non-group attributes
	if attr.Key == "" {
		return event
	}

	key := joinPrefix(prefix, attr.Key)
	val := attr.Value

	switch val.Kind() {
	case slog.KindString:
		event = event.Str(key, val.String())
	case slog.KindInt64:
		event = event.Int64(key, val.Int64())
	case slog.KindUint64:
		event = event.Uint64(key, val.Uint64())
	case slog.KindFloat64:
		event = event.Float64(key, val.Float64())
	case slog.KindBool:
		event = event.Bool(key, val.Bool())
	case slog.KindDuration:
		event = event.Dur(key, val.Duration())
	case slog.KindTime:
		event = event.Time(key, val.Time())
	case slog.KindAny:
		v := val.Any()
		switch cv := v.(type) {
		case error:
			event = event.AnErr(key, cv)
		case time.Duration:
			event = event.Dur(key, cv)
		case time.Time:
			event = event.Time(key, cv)
		case []byte:
			event = event.Bytes(key, cv)
		default:
			event = event.Interface(key, v)
		}
	default:
		event = event.Interface(key, val.Any())
	}

	return event
}

// Verify at compile time that loggerHandler satisfies the slog.Handler interface.
var _ slog.Handler = (*loggerHandler)(nil)
