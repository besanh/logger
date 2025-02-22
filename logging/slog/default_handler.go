package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
)

type defaultHandler struct {
	slog.Handler
}

func NewDefaultHandler(w io.Writer, opts *slog.HandlerOptions) *defaultHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &defaultHandler{
		slog.NewJSONHandler(w, opts),
	}
}

func (h defaultHandler) Handle(ctx context.Context, r slog.Record) error {
	if logId, ok := ctx.Value("trace_id").(string); ok {
		r.Add("trace_id", logId)
	} else {
		r.Add("trace_id", "unknown")
	}
	_, path, numLine, _ := runtime.Caller(4)
	srcFile := filepath.Base(path)
	r.Add(slog.String("file", fmt.Sprintf("%s:%d", srcFile, numLine)))
	return h.Handler.Handle(ctx, r)
}
