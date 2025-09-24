// Package logger - пакет для работы с логгером
package logger

import (
	"os"

	"github.com/rlapenok/toolbox/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogFormat - format of the log
type LogFormat string

const (
	FormatLogfmt LogFormat = "logfmt"
	FormatJSON   LogFormat = "json"
)

// New - create new logger
func New(cfg Config) (*zap.Logger, error) {
	// create level
	level, err := zapcore.ParseLevel(cfg.GetLevel())
	if err != nil {
		message := err.Error()
		return nil, errors.New(errors.InvalidParameter, message)
	}

	// create encoder
	encoder := makeEncoder(string(cfg.GetFormat()))

	stdout := zapcore.Lock(zapcore.AddSync(os.Stdout))
	core := zapcore.NewCore(encoder, stdout, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= level
	}))

	log := zap.New(core)

	// meta
	meta := cfg.GetMeta()
	if len(meta) > 0 {
		fields := make([]zap.Field, 0, len(meta))
		for k, v := range meta {
			fields = append(fields, zap.Any(k, v))
		}
		log = log.With(fields...)
	}

	return log, nil
}

func makeEncoder(format string) zapcore.Encoder {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	switch format {
	case "json":
		return zapcore.NewJSONEncoder(encCfg)
	default:
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encCfg)
	}
}
