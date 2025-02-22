package slog

import "context"

type (
	ILogger interface {
		Info(args ...any)
		Warn(args ...any)
		Error(args ...any)
		Debug(args ...any)
	}

	IFormatLogger interface {
		Infof(msg string, args ...any)
		Warnf(msg string, args ...any)
		Errorf(msg string, args ...any)
		Debugf(msg string, args ...any)
	}

	IContextLogger interface {
		InfoContext(ctx context.Context, args ...any)
		WarnContext(ctx context.Context, args ...any)
		ErrorContext(ctx context.Context, args ...any)
		DebugContext(ctx context.Context, args ...any)
	}

	IContextFormatLogger interface {
		InfofContext(ctx context.Context, msg string, args ...any)
		WarnfContext(ctx context.Context, msg string, args ...any)
		ErrorfContext(ctx context.Context, msg string, args ...any)
		DebugfContext(ctx context.Context, msg string, args ...any)
	}

	IFullLogger interface {
		ILogger
		IFormatLogger
		IContextLogger
		IContextFormatLogger
	}
)

type (
	Control interface {
		SetLevel(level Level)
	}
)
