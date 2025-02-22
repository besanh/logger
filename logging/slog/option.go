package slog

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/gookit/slog/rotatefile"
)

type (
	Level     slog.Level
	Formatter string

	coreConfig struct {
		level     *slog.LevelVar
		formatter Formatter

		opt                *slog.HandlerOptions
		writer             io.Writer
		withLevel          bool
		withHandlerOptions bool

		// set options
		isWithTraceId bool

		isWithFileSource bool

		// optional: fetch attributes from context
		AddSource   bool
		ReplaceAttr func(groups []string, a slog.Attr) slog.Attr

		// optional: connection to Fluentd
		isUseFluent  bool
		FluentClient *fluent.Fluent // fluent client if isUseFluent is true
		Tag          string

		// optional: customize json payload builder
		Converter Converter

		attrs           []slog.Attr
		AttrFromContext []func(ctx context.Context) []slog.Attr
	}

	config struct {
		coreConfig coreConfig
	}
)

const (
	LEVEL_DEBUG Level = -4
	LEVEL_INFO  Level = 0
	LEVEL_WARN  Level = 4
	LEVEL_ERROR Level = 8
	LEVEL_FATAL Level = 12
	LEVEL_TRACE Level = 16

	FORMAT_JSON Formatter = "json"
	FORMAT_TEXT Formatter = "text"
)

// Option slog option
type Option interface {
	apply(cfg *config)
}

type option func(cfg *config)

func (fn option) apply(cfg *config) {
	fn(cfg)
}

// default config
func defaultConfig() *config {
	coreConfig := defaultCoreConfig()
	return &config{
		coreConfig: *coreConfig,
	}
}

// default core config
func defaultCoreConfig() *coreConfig {
	level := new(slog.LevelVar)
	level.Set(slog.LevelInfo)
	return &coreConfig{
		opt: &slog.HandlerOptions{
			Level: level,
		},
		writer:             os.Stdout,
		level:              level,
		withLevel:          false,
		withHandlerOptions: false,
		formatter:          FORMAT_JSON,
		isUseFluent:        false,
		AddSource:          false,
		isWithTraceId:      false,
		isWithFileSource:   false,
		attrs:              []slog.Attr{},
	}
}

// WithHandlerOptions slog handler-options
func WithHandlerOptions(opt *slog.HandlerOptions) Option {
	return option(func(cfg *config) {
		cfg.coreConfig.opt = opt
		cfg.coreConfig.withHandlerOptions = true
	})
}

// WithOutput slog writer
func WithOutput(iow io.Writer) Option {
	return option(func(cfg *config) {
		cfg.coreConfig.writer = iow
	})
}

// WithLevel slog level
func WithLevel(level Level) Option {
	return option(func(cfg *config) {
		cfg.coreConfig.level.Set(tranSLevel(level))
		cfg.coreConfig.withLevel = true
	})
}

// WithRotateFile rotate file option
func WithRotateFile(f string) Option {
	rotateWriter, err := rotatefile.NewConfig(f).Create()
	if err != nil {
		panic(err)
	}
	w := io.MultiWriter(os.Stdout, rotateWriter)
	return option(func(cfg *config) {
		cfg.coreConfig.writer = w
	})
}

// WithFormatter formatter
// default json.
// Enum: json or text
func WithFormatter(formatter Formatter) Option {
	return option(func(cfg *config) {
		cfg.coreConfig.formatter = formatter
	})
}

func WithTraceId() Option {
	return option(func(cfg *config) {
		cfg.coreConfig.isWithTraceId = true
	})
}

func WithFileSource() Option {
	return option(func(cfg *config) {
		cfg.coreConfig.isWithFileSource = true
	})
}

// WithFluentd fluentd config
func WithFluentd(client *fluent.Fluent, tag string) Option {
	return option(func(cfg *config) {
		cfg.coreConfig.isUseFluent = true
		cfg.coreConfig.FluentClient = client
		cfg.coreConfig.Tag = tag
		cfg.coreConfig.AttrFromContext = []func(ctx context.Context) []slog.Attr{}
	})
}

func WithAttrs(attrs ...slog.Attr) Option {
	return option(func(cfg *config) {
		cfg.coreConfig.attrs = attrs
	})
}
