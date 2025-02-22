package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

// TraceLogger is a logger that supports trace with opentelemetry trace
type TraceLogger struct {
	l      *slog.Logger
	config *config
}

// NewTraceLogger new trace logger
func NewTraceLogger(opts ...Option) *TraceLogger {
	config := defaultConfig()
	for _, opt := range opts {
		opt.apply(config)
	}

	// When user set the handlerOptions level but not set with coreconfig level
	if !config.coreConfig.withLevel && config.coreConfig.withHandlerOptions && config.coreConfig.opt.Level != nil {
		lvl := &slog.LevelVar{}
		lvl.Set(config.coreConfig.opt.Level.Level())
		config.coreConfig.level = lvl
	}
	config.coreConfig.opt.Level = config.coreConfig.level

	logger := slog.New(NewTraceHandler(config.coreConfig.writer, config.coreConfig.opt, config.traceConfig))

	return &TraceLogger{
		l:      logger,
		config: config,
	}
}

func (l *TraceLogger) Log(level Level, msg string) {
	logger := l.l.With()
	logger.Log(context.TODO(), tranSLevel(level), msg)
}

func (l *TraceLogger) Logf(level Level, format string, kvs ...interface{}) {
	logger := l.l.With()
	msg := getMessage(format, kvs)
	logger.Log(context.TODO(), tranSLevel(level), msg)
}

func (l *TraceLogger) LogCtxf(level Level, ctx context.Context, format string, kvs ...interface{}) {
	logger := l.l.With()
	msg := getMessage(format, kvs)
	logger.Log(ctx, tranSLevel(level), msg)
}

func (l *TraceLogger) Debug(args ...any) {
	l.Log(Level(LEVEL_DEBUG), fmt.Sprint(args...))
}

func (l *TraceLogger) Info(args ...any) {
	l.Log(Level(LEVEL_INFO), fmt.Sprint(args...))
}

func (l *TraceLogger) Warn(args ...any) {
	l.Log(Level(LEVEL_WARN), fmt.Sprint(args...))
}

func (l *TraceLogger) Error(args ...any) {
	l.Log(Level(LEVEL_ERROR), fmt.Sprint(args...))
}

func (l *TraceLogger) Debugf(msg string, args ...any) {
	l.Logf(Level(LEVEL_DEBUG), msg, args...)
}

func (l *TraceLogger) Infof(msg string, args ...any) {
	l.Logf(Level(LEVEL_INFO), msg, args...)
}

func (l *TraceLogger) Warnf(msg string, args ...any) {
	l.Logf(Level(LEVEL_WARN), msg, args...)
}

func (l *TraceLogger) Errorf(msg string, args ...any) {
	l.Logf(Level(LEVEL_ERROR), msg, args...)
}

func (l *TraceLogger) DebugContext(ctx context.Context, args ...any) {
	l.LogCtxf(Level(LEVEL_DEBUG), ctx, fmt.Sprint(args...))
}

func (l *TraceLogger) InfoContext(ctx context.Context, args ...any) {
	l.LogCtxf(Level(LEVEL_INFO), ctx, fmt.Sprint(args...))
}

func (l *TraceLogger) WarnContext(ctx context.Context, args ...any) {
	l.LogCtxf(Level(LEVEL_WARN), ctx, fmt.Sprint(args...))
}

func (l *TraceLogger) ErrorContext(ctx context.Context, args ...any) {
	l.LogCtxf(Level(LEVEL_ERROR), ctx, fmt.Sprint(args...))
}

func (l *TraceLogger) DebugfContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(Level(LEVEL_DEBUG), ctx, msg, args...)
}

func (l *TraceLogger) InfofContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(Level(LEVEL_INFO), ctx, msg, args...)
}

func (l *TraceLogger) WarnfContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(Level(LEVEL_WARN), ctx, msg, args...)
}

func (l *TraceLogger) ErrorfContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(Level(LEVEL_ERROR), ctx, msg, args...)
}

func (l *TraceLogger) SetLevel(level Level) {
	lvl := tranSLevel(level)
	l.config.coreConfig.level.Set(lvl)
}

func (l *TraceLogger) SetOutput(writer io.Writer) {
	log := slog.New(NewDefaultHandler(writer, l.config.coreConfig.opt))
	l.config.coreConfig.writer = writer
	l.l = log
}
