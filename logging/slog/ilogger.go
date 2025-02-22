package slog

import (
	"context"
	"log/slog"
)

var logger IFullLogger = NewSLogger()

// SetLogger sets the default logger.
// Note that this method is not concurrent-safe and must not be called
// after the use of DefaultLogger and global functions in this package.
func SetLogger(l IFullLogger) {
	logger = l
}

// DefaultLogger return the default logger.
func GetLogger() IFullLogger {
	return logger
}

// Info calls the default logger's Info method.
func Info(args ...any) {
	logger.Info(args...)
}

// Warn calls the default logger's Warn method.
func Warn(args ...any) {
	logger.Warn(args...)
}

// Error calls the default logger's Error method.
func Error(args ...any) {
	logger.Error(args...)
}

// Debug calls the default logger's Debug method.
func Debug(args ...any) {
	logger.Debug(args...)
}

// Infof calls the default logger's Infof method.
func Infof(msg string, args ...any) {
	logger.Infof(msg, args...)
}

// Warnf calls the default logger's Warnf method.
func Warnf(msg string, args ...any) {
	logger.Warnf(msg, args...)
}

// Errorf calls the default logger's Errorf method.
func Errorf(msg string, args ...any) {
	logger.Errorf(msg, args...)
}

// Debugf calls the default logger's Debugf method.
func Debugf(msg string, args ...any) {
	logger.Debugf(msg, args...)
}

// InfoContext calls the default logger's InfoContext method.
func InfoContext(ctx context.Context, args ...any) {
	logger.InfoContext(ctx, args...)
}

// WarnContext calls the default logger's WarnContext method.
func WarnContext(ctx context.Context, args ...any) {
	logger.WarnContext(ctx, args...)
}

// ErrorContext calls the default logger's ErrorContext method.
func ErrorContext(ctx context.Context, args ...any) {
	logger.ErrorContext(ctx, args...)
}

// DebugContext calls the default logger's DebugContext method.
func DebugContext(ctx context.Context, args ...any) {
	logger.DebugContext(ctx, args...)
}

// InfofContext calls the default logger's InfofContext method.
func InfofContext(ctx context.Context, msg string, args ...any) {
	logger.InfofContext(ctx, msg, args...)
}

// WarnfContext calls the default logger's WarnfContext method.
func WarnfContext(ctx context.Context, msg string, args ...any) {
	logger.WarnfContext(ctx, msg, args...)
}

// ErrorfContext calls the default logger's ErrorfContext method.
func ErrorfContext(ctx context.Context, msg string, args ...any) {
	logger.ErrorfContext(ctx, msg, args...)
}

// DebugfContext calls the default logger's DebugfContext method.
func DebugfContext(ctx context.Context, msg string, args ...any) {
	logger.DebugfContext(ctx, msg, args...)
}

func LogWithAttrs(level Level, ctx context.Context, msg string, attrs ...slog.Attr) {
	logger.LogWithAttrs(level, ctx, msg, attrs...)
}
