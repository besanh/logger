package slog

import (
	"log/slog"
	"time"

	slogcommon "github.com/samber/slog-common"
)

var SourceKey = "source"

type Converter func(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record, tag string) map[string]string

func DefaultConverter(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record, tag string) map[string]string {
	// aggregate all attributes
	attrs := slogcommon.AppendRecordAttrsToAttrs(loggerAttr, groups, record)

	// developer formatters
	if addSource {
		attrs = append(attrs, slogcommon.Source(SourceKey, record))
	}
	attrs = slogcommon.ReplaceAttrs(replaceAttr, []string{}, attrs...)
	attrs = slogcommon.RemoveEmptyAttrs(attrs)

	// handler formatter
	log := map[string]string{
		"timestamp": record.Time.UTC().Format(time.RFC3339),
		"level":     record.Level.String(),
		"message":   record.Message,
		"tag":       tag,
	}

	extra := AttrsToString(attrs...)

	for k, v := range extra {
		if _, ok := log[k]; ok {
			continue
		}
		log[k] = v
	}

	return log
}
