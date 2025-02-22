package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"

	slogcommon "github.com/samber/slog-common"
)

type defaultHandler struct {
	slog.Handler
	config *coreConfig
	attrs  []slog.Attr
	groups []string
}

func NewDefaultHandler(w io.Writer, opts *slog.HandlerOptions, config *coreConfig) *defaultHandler {
	// if opts == nil {
	// 	opts = &slog.HandlerOptions{}
	// }

	if config.isUseFluent {
		// check config is valid
		if config.FluentClient == nil {
			panic("missing Fluent client")
		}

		if config.Converter == nil {
			config.Converter = DefaultConverter
		}

		if config.AttrFromContext == nil {
			config.AttrFromContext = []func(ctx context.Context) []slog.Attr{}
		}
	}
	var handler slog.Handler
	switch config.formatter {
	case FORMAT_TEXT:
		handler = slog.NewTextHandler(w, opts)
	default:
		handler = slog.NewJSONHandler(w, opts)
	}
	return &defaultHandler{
		Handler: handler,
		config:  config,
	}
}

var _ slog.Handler = (*defaultHandler)(nil)

func (h *defaultHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.config.level.Level()
}

func (h *defaultHandler) Handle(ctx context.Context, record slog.Record) error {
	if h.config.isWithTraceId {
		if logId, ok := ctx.Value(keyTraceId).(string); ok {
			record.Add("trace_id", logId)
		} else {
			record.Add("trace_id", "unknown")
		}
	}

	if len(h.config.attrs) > 0 {
		record.AddAttrs(h.config.attrs...)
	}
	// developer formatters
	if h.config.isWithFileSource {
		_, path, numLine, _ := runtime.Caller(6)
		srcFile := filepath.Base(path)
		record.Add(slog.String("file", fmt.Sprintf("%s:%d", srcFile, numLine)))
	}

	if h.config.isUseFluent {
		go h.postToFluent(ctx, record)
	}
	// handler formatter
	return h.Handler.Handle(ctx, record)
}

func (h *defaultHandler) postToFluent(ctx context.Context, record slog.Record) error {
	tag := h.getTag(&record)
	fromContext := ContextExtractor(ctx, h.config.AttrFromContext)
	message := h.config.Converter(h.config.AddSource, h.config.ReplaceAttr, append(h.attrs, fromContext...), h.groups, &record, tag)

	return h.config.FluentClient.PostWithTime(tag, record.Time, message)
}

func (h *defaultHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &defaultHandler{
		config: h.config,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *defaultHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	return &defaultHandler{
		config: h.config,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}

func (h *defaultHandler) getTag(record *slog.Record) string {
	tag := h.config.Tag

	for i := range h.attrs {
		if h.attrs[i].Key == "tag" && h.attrs[i].Value.Kind() == slog.KindString {
			tag = h.attrs[i].Value.String()
			break
		}
	}

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "tag" && attr.Value.Kind() == slog.KindString {
			tag = attr.Value.String()
			return false
		}
		return true
	})

	return tag
}
