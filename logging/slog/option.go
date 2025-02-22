package slog

import (
	"io"
	"log/slog"
	"os"

	"github.com/gookit/slog/rotatefile"
)

type (
	Level int

	coreConfig struct {
		opt                *slog.HandlerOptions
		writer             io.Writer
		level              *slog.LevelVar
		withLevel          bool
		withHandlerOptions bool
	}

	config struct {
		coreConfig  coreConfig
		traceConfig *traceConfig
	}
)

const (
	LEVEL_DEBUG int = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
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
		traceConfig: &traceConfig{
			recordStackTraceInSpan: true,
			errorSpanLevel:         slog.LevelError,
		},
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
func WithLevel(lvl *slog.LevelVar) Option {
	return option(func(cfg *config) {
		cfg.coreConfig.level = lvl
		cfg.coreConfig.withLevel = true
	})
}

// WithTraceErrorSpanLevel trace error span level option
func WithTraceErrorSpanLevel(level slog.Level) Option {
	return option(func(cfg *config) {
		cfg.traceConfig.errorSpanLevel = level
	})
}

// WithRecordStackTraceInSpan record stack track option
func WithRecordStackTraceInSpan(recordStackTraceInSpan bool) Option {
	return option(func(cfg *config) {
		cfg.traceConfig.recordStackTraceInSpan = recordStackTraceInSpan
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
