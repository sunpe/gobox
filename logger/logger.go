package logger

import (
	"context"
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/slog"
)

// Init logger
func Init(opts ...Option) {
	o := defaultOption
	for _, opt := range opts {
		opt(&o)
	}
	slog.SetDefault(o.newLogger())
}

type Option func(option *option)

// WithWriter set writer for logger. default is os.Stdout
func WithWriter(writer io.Writer) Option {
	return func(o *option) {
		o.writer = writer
	}
}

// WithLevel set level for logger. default is LevelInfo
func WithLevel(level LogLevel) Option {
	return func(o *option) {
		o.level = level
	}
}

// JSONOutput set whether output json format. default is false
func JSONOutput() Option {
	return func(o *option) {
		o.json = true
	}
}

// TextOutput set whether output text format. default is false
func TextOutput() Option {
	return func(o *option) {
		o.text = true
	}
}

// WithAttr set attributes for logger. default is empty
func WithAttr(key string, value any) Option {
	return func(o *option) {
		o.attrs[key] = value
	}
}

func WithSource() Option {
	return func(o *option) {
		o.addSource = true
	}
}

// LogLevel is the level of a logger.
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
)

// Debug show debug log
func Debug(msg string, args ...any) {
	DebugWithCtx(context.Background(), msg, args...)
}

func DebugWithCtx(ctx context.Context, msg string, args ...any) {
	slog.DebugCtx(ctx, msg, args...)
}

func DebugF(format string, v ...any) {
	DebugFWithCtx(context.Background(), format, v...)
}

func DebugFWithCtx(ctx context.Context, format string, v ...any) {
	slog.DebugCtx(ctx, fmt.Sprintf(format, v...))
}

func Info(msg string, args ...any) {
	InfoWithCtx(context.Background(), msg, args...)
}

func InfoWithCtx(ctx context.Context, msg string, args ...any) {
	slog.InfoCtx(ctx, msg, args...)
}

func InfoF(format string, v ...any) {
	InfoFWithCtx(context.Background(), format, v...)
}

func InfoFWithCtx(ctx context.Context, format string, v ...any) {
	slog.InfoCtx(ctx, fmt.Sprintf(format, v...))
}

func Warn(msg string, args ...any) {
	WarnWithCtx(context.Background(), msg, args...)
}

func WarnWithCtx(ctx context.Context, msg string, args ...any) {
	slog.WarnCtx(ctx, msg, args...)
}

func WarnF(format string, v ...any) {
	WarnFWithCtx(context.Background(), format, v...)
}

func WarnFWithCtx(ctx context.Context, format string, v ...any) {
	slog.WarnCtx(ctx, fmt.Sprintf(format, v...))
}

func Error(msg string, args ...any) {
	ErrorWithCtx(context.Background(), msg, args...)
}

func ErrorWithCtx(ctx context.Context, msg string, args ...any) {
	slog.ErrorCtx(ctx, msg, args...)
}

func ErrorF(format string, v ...any) {
	ErrorFWithCtx(context.Background(), format, v...)
}

func ErrorFWithCtx(ctx context.Context, format string, v ...any) {
	slog.ErrorCtx(ctx, fmt.Sprintf(format, v...))
}

func Panic(msg string, args ...any) {
	PanicWithCtx(context.Background(), msg, args...)
}

func PanicWithCtx(ctx context.Context, msg string, args ...any) {
	slog.Log(ctx, slogLevelPanic, msg, args...)
	messages := make([]interface{}, 0, len(args)+1)
	messages = append(messages, msg)
	messages = append(messages, args...)
	panic(fmt.Sprint(messages...))
}

func PanicF(format string, v ...any) {
	PanicFWithCtx(context.Background(), format, v...)
}

func PanicFWithCtx(ctx context.Context, format string, v ...any) {
	slog.Log(ctx, slogLevelPanic, fmt.Sprintf(format, v...))
	panic(fmt.Sprintf(format, v...))
}

func init() {
	Init()
}

var defaultOption = option{
	writer:    os.Stdout,
	addSource: false,
	level:     LevelInfo,
	json:      false,
	attrs:     map[string]any{},
}

type option struct {
	writer    io.Writer
	addSource bool
	level     LogLevel
	json      bool
	text      bool
	attrs     map[string]any
}

var levelMap = map[LogLevel]slog.Level{
	LevelDebug: slog.LevelDebug,
	LevelInfo:  slog.LevelInfo,
	LevelWarn:  slog.LevelWarn,
	LevelError: slog.LevelError,
	LevelPanic: slogLevelPanic,
}

const slogLevelPanic = slog.Level(12)

var sLogLevelName = map[slog.Level]string{
	slogLevelPanic: "PANIC",
}

func (o *option) newLogger() *slog.Logger {
	var h slog.Handler
	handlerOps := slog.HandlerOptions{
		AddSource: o.addSource,
		Level:     levelMap[o.level],
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				name, ok := sLogLevelName[level]
				if !ok {
					name = level.String()
				}
				a.Value = slog.StringValue(name)
			}
			return a
		},
	}

	if o.json {
		h = slog.NewJSONHandler(o.writer, &handlerOps)
	} else {
		h = slog.NewTextHandler(o.writer, &handlerOps)
	}
	if len(o.attrs) > 0 {
		attrs := make([]slog.Attr, 0, len(o.attrs))
		for k, v := range o.attrs {
			attrs = append(attrs, slog.Attr{
				Key:   k,
				Value: slog.AnyValue(v),
			})
		}
		h = h.WithAttrs(attrs)
	}

	return slog.New(h)
}
