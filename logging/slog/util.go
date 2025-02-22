package slog

import (
	"context"
	"encoding"
	"fmt"
	"log/slog"
	"strconv"
)

// get format msg
func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

// Adapt level to slog level
func tranSLevel(level Level) (lvl slog.Level) {
	switch level {
	case LEVEL_DEBUG:
		lvl = slog.LevelDebug
	case LEVEL_INFO:
		lvl = slog.LevelInfo
	case LEVEL_WARN:
		lvl = slog.LevelWarn
	case LEVEL_ERROR:
		lvl = slog.LevelError
	default:
		lvl = slog.LevelDebug
	}
	return
}

type ContextKey string

const keyTraceId = ContextKey("trace_id")

func SetContextTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, keyTraceId, traceId)
}

func AttrsToString(attrs ...slog.Attr) map[string]string {
	output := map[string]string{}

	for i := range attrs {
		attr := attrs[i]
		k, v := attr.Key, attr.Value
		output[k] = ValueToString(v)
	}

	return output
}

func ValueToString(v slog.Value) string {
	switch v.Kind() {
	case slog.KindAny:
		return AnyValueToString(v)
	case slog.KindLogValuer:
		return AnyValueToString(v)
	case slog.KindGroup:
		return AnyValueToString(v)
	case slog.KindInt64:
		return fmt.Sprintf("%d", v.Int64())
	case slog.KindUint64:
		return fmt.Sprintf("%d", v.Uint64())
	case slog.KindFloat64:
		return fmt.Sprintf("%f", v.Float64())
	case slog.KindString:
		return v.String()
	case slog.KindBool:
		return strconv.FormatBool(v.Bool())
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().UTC().String()
	default:
		return AnyValueToString(v)
	}
}

func AnyValueToString(v slog.Value) string {
	if tm, ok := v.Any().(encoding.TextMarshaler); ok {
		data, err := tm.MarshalText()
		if err != nil {
			return ""
		}
		return string(data)
	}

	return fmt.Sprintf("%+v", v.Any())
}

func ContextExtractor(ctx context.Context, fns []func(ctx context.Context) []slog.Attr) []slog.Attr {
	attrs := []slog.Attr{}
	for _, fn := range fns {
		attrs = append(attrs, fn(ctx)...)
	}
	return attrs
}

func ExtractFromContext(keys ...any) func(ctx context.Context) []slog.Attr {
	return func(ctx context.Context) []slog.Attr {
		attrs := []slog.Attr{}
		for _, key := range keys {
			attrs = append(attrs, slog.Any(key.(string), ctx.Value(key)))
		}
		return attrs
	}
}
