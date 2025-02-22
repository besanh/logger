package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type DefaultLogger struct {
	l      *slog.Logger
	config *config
}

func NewDefaultLogger(opts ...Option) *DefaultLogger {
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
	logger := slog.New(NewDefaultHandler(config.coreConfig.writer, config.coreConfig.opt))
	return &DefaultLogger{
		l:      logger,
		config: config,
	}
}

func (l *DefaultLogger) Log(level Level, msg string) {
	logger := l.l.With()
	logger.Log(context.TODO(), tranSLevel(level), msg)
}

func (l *DefaultLogger) Logf(level Level, format string, kvs ...interface{}) {
	logger := l.l.With()
	msg := getMessage(format, kvs)
	logger.Log(context.TODO(), tranSLevel(level), msg)
}

func (l *DefaultLogger) LogCtxf(level Level, ctx context.Context, format string, kvs ...interface{}) {
	logger := l.l.With()
	msg := getMessage(format, kvs)
	logger.Log(ctx, tranSLevel(level), msg)
}

func (l *DefaultLogger) Debug(args ...any) {
	l.Log(Level(LEVEL_DEBUG), fmt.Sprint(args...))
}

func (l *DefaultLogger) Info(args ...any) {
	l.Log(Level(LEVEL_INFO), fmt.Sprint(args...))
}

func (l *DefaultLogger) Warn(args ...any) {
	l.Log(Level(LEVEL_WARN), fmt.Sprint(args...))
}

func (l *DefaultLogger) Error(args ...any) {
	l.Log(Level(LEVEL_ERROR), fmt.Sprint(args...))
}

func (l *DefaultLogger) Debugf(msg string, args ...any) {
	l.Logf(Level(LEVEL_DEBUG), msg, args...)
}

func (l *DefaultLogger) Infof(msg string, args ...any) {
	l.Logf(Level(LEVEL_INFO), msg, args...)
}

func (l *DefaultLogger) Warnf(msg string, args ...any) {
	l.Logf(Level(LEVEL_WARN), msg, args...)
}

func (l *DefaultLogger) Errorf(msg string, args ...any) {
	l.Logf(Level(LEVEL_ERROR), msg, args...)
}

func (l *DefaultLogger) DebugContext(ctx context.Context, args ...any) {
	l.LogCtxf(Level(LEVEL_DEBUG), ctx, fmt.Sprint(args...))
}

func (l *DefaultLogger) InfoContext(ctx context.Context, args ...any) {
	l.LogCtxf(Level(LEVEL_INFO), ctx, fmt.Sprint(args...))
}

func (l *DefaultLogger) WarnContext(ctx context.Context, args ...any) {
	l.LogCtxf(Level(LEVEL_WARN), ctx, fmt.Sprint(args...))
}

func (l *DefaultLogger) ErrorContext(ctx context.Context, args ...any) {
	l.LogCtxf(Level(LEVEL_ERROR), ctx, fmt.Sprint(args...))
}

func (l *DefaultLogger) DebugfContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(Level(LEVEL_DEBUG), ctx, msg, args...)
}

func (l *DefaultLogger) InfofContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(Level(LEVEL_INFO), ctx, msg, args...)
}

func (l *DefaultLogger) WarnfContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(Level(LEVEL_WARN), ctx, msg, args...)
}

func (l *DefaultLogger) ErrorfContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(Level(LEVEL_ERROR), ctx, msg, args...)
}

func (l *DefaultLogger) SetLevel(level Level) {
	lvl := tranSLevel(level)
	l.config.coreConfig.level.Set(lvl)
}

func (l *DefaultLogger) SetOutput(writer io.Writer) {
	log := slog.New(NewDefaultHandler(writer, l.config.coreConfig.opt))
	l.config.coreConfig.writer = writer
	l.l = log
}
