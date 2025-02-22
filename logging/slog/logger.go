package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type SLogger struct {
	l      *slog.Logger
	config *config
}

func NewSLogger(opts ...Option) *SLogger {
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
	logger := slog.New(NewDefaultHandler(config.coreConfig.writer, config.coreConfig.opt, &config.coreConfig))
	return &SLogger{
		l:      logger,
		config: config,
	}
}

var _ ILogger = (*SLogger)(nil)

func (l *SLogger) Log(level Level, msg string) {
	logger := l.l.With()
	logger.Log(context.TODO(), tranSLevel(level), msg)
}

func (l *SLogger) Logf(level Level, format string, kvs ...interface{}) {
	logger := l.l.With()
	msg := getMessage(format, kvs)
	logger.Log(context.TODO(), tranSLevel(level), msg)
}

func (l *SLogger) LogCtxf(level Level, ctx context.Context, format string, kvs ...interface{}) {
	logger := l.l.With()
	msg := getMessage(format, kvs)
	logger.Log(ctx, tranSLevel(level), msg)
}

func (l *SLogger) LogWithAttrs(level Level, ctx context.Context, msg string, attrs ...slog.Attr) {
	logger := l.l.With()
	logger.LogAttrs(ctx, tranSLevel(level), msg, attrs...)
}

func (l *SLogger) Debug(args ...any) {
	l.Log(LEVEL_DEBUG, fmt.Sprint(args...))
}

func (l *SLogger) Info(args ...any) {
	l.Log(LEVEL_INFO, fmt.Sprint(args...))
}

func (l *SLogger) Warn(args ...any) {
	l.Log(LEVEL_WARN, fmt.Sprint(args...))
}

func (l *SLogger) Error(args ...any) {
	l.Log(LEVEL_ERROR, fmt.Sprint(args...))
}

func (l *SLogger) Debugf(msg string, args ...any) {
	l.Logf(LEVEL_DEBUG, msg, args...)
}

func (l *SLogger) Infof(msg string, args ...any) {
	l.Logf(LEVEL_INFO, msg, args...)
}

func (l *SLogger) Warnf(msg string, args ...any) {
	l.Logf(LEVEL_WARN, msg, args...)
}

func (l *SLogger) Errorf(msg string, args ...any) {
	l.Logf(LEVEL_ERROR, msg, args...)
}

func (l *SLogger) DebugContext(ctx context.Context, args ...any) {
	l.LogCtxf(LEVEL_DEBUG, ctx, fmt.Sprint(args...))
}

func (l *SLogger) InfoContext(ctx context.Context, args ...any) {
	l.LogCtxf(LEVEL_INFO, ctx, fmt.Sprint(args...))
}

func (l *SLogger) WarnContext(ctx context.Context, args ...any) {
	l.LogCtxf(LEVEL_WARN, ctx, fmt.Sprint(args...))
}

func (l *SLogger) ErrorContext(ctx context.Context, args ...any) {
	l.LogCtxf(LEVEL_ERROR, ctx, fmt.Sprint(args...))
}

func (l *SLogger) DebugfContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(LEVEL_DEBUG, ctx, msg, args...)
}

func (l *SLogger) InfofContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(LEVEL_INFO, ctx, msg, args...)
}

func (l *SLogger) WarnfContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(LEVEL_WARN, ctx, msg, args...)
}

func (l *SLogger) ErrorfContext(ctx context.Context, msg string, args ...any) {
	l.LogCtxf(LEVEL_ERROR, ctx, msg, args...)
}

func (l *SLogger) SetLevel(level Level) {
	lvl := tranSLevel(level)
	l.config.coreConfig.level.Set(lvl)
}

func (l *SLogger) SetOutput(writer io.Writer) {
	log := slog.New(NewDefaultHandler(writer, l.config.coreConfig.opt, &l.config.coreConfig))
	l.config.coreConfig.writer = writer
	l.l = log
}
