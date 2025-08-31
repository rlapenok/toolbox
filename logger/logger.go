// Package logger - пакет для работы с логгером
package logger

import (
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ConsoleMode - режим вывода в консоль
type (
	// ConsoleMode - режим вывода в консоль
	ConsoleMode string

	// LogFormat - формат вывода в консоль
	LogFormat string
)

const (
	// ConsoleAll - всё в stdout (включая ошибки)
	ConsoleAll ConsoleMode = "all"

	// ConsoleSplit - error+ -> stderr, остальное -> stdout
	ConsoleSplit ConsoleMode = "split"

	// FormatLogfmt - формат вывода в консоль
	FormatLogfmt LogFormat = "logfmt"

	// FormatJSON - формат вывода в консоль
	FormatJSON LogFormat = "json"
)

type Config struct {
	// Консоль
	ConsoleEnabled bool
	ConsoleFormat  LogFormat
	ConsoleMode    ConsoleMode
	ConsoleLevel   string // "debug"|"info"|...

	// Файл
	FileEnabled    bool
	FilePath       string
	FileFormat     LogFormat // "logfmt"|"json" (default logfmt)
	FileLevel      string    // "debug"|"info"|...
	FileMaxSizeMB  int
	FileMaxBackups int
	FileMaxAgeDays int
	FileCompress   bool

	Meta map[string]any // метаданные для всех логов
}

// New - создание нового логгера
func New(cfg Config) (*zap.Logger, error) {
	var consoleLvl zapcore.Level
	var fileLvl zapcore.Level

	// Проверяем уровень консоли
	if cfg.ConsoleEnabled {
		if l, err := zapcore.ParseLevel(cfg.ConsoleLevel); err != nil {
			return nil, fmt.Errorf("logger: bad console level: %w", err)
		} else {
			consoleLvl = l
		}
	}

	// Проверяем уровень файла
	if cfg.FileEnabled {
		if l, err := zapcore.ParseLevel(cfg.FileLevel); err != nil {
			return nil, fmt.Errorf("logger: bad file min level: %w", err)
		} else {
			fileLvl = l
		}
	}

	// Создаем слайс для хранения ядер
	var cores []zapcore.Core

	// Файл
	if cfg.FileEnabled {
		if cfg.FilePath == "" {
			return nil, errors.New("logger: file output enabled but FilePath is empty")
		}
		fileEnc := makeEncoder(string(cfg.FileFormat))
		lj := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    max(1, cfg.FileMaxSizeMB),
			MaxBackups: max(0, cfg.FileMaxBackups),
			MaxAge:     max(0, cfg.FileMaxAgeDays),
			Compress:   cfg.FileCompress,
		}
		fileWS := zapcore.Lock(zapcore.AddSync(lj))
		cores = append(cores, zapcore.NewCore(fileEnc, fileWS, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= fileLvl
		})))
	}

	// Консоль
	if cfg.ConsoleEnabled {
		consoleEnc := makeEncoder(string(cfg.ConsoleFormat))

		switch cfg.ConsoleMode {
		// ConsoleSplit - ошибки в stderr, остальное в stdout
		case ConsoleSplit:
			// non-errors -> stdout
			stdout := zapcore.Lock(zapcore.AddSync(os.Stdout))
			cores = append(cores, zapcore.NewCore(consoleEnc, stdout, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
				return l >= consoleLvl && l < zapcore.ErrorLevel
			})))
			// errors+ -> stderr
			stderr := zapcore.Lock(zapcore.AddSync(os.Stderr))
			cores = append(cores, zapcore.NewCore(consoleEnc, stderr, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
				return l >= zapcore.ErrorLevel
			})))
		// ConsoleAll - всё в stdout
		default:
			stdout := zapcore.Lock(zapcore.AddSync(os.Stdout))
			cores = append(cores, zapcore.NewCore(consoleEnc, stdout, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
				return l >= consoleLvl
			})))
		}
	}

	core := zapcore.NewTee(cores...)

	// сглаживание бурстов (по желанию)
	core = zapcore.NewSamplerWithOptions(core, time.Second, 100, 100)

	log := zap.New(core)

	// метаданные
	if len(cfg.Meta) > 0 {
		fields := make([]zap.Field, 0, len(cfg.Meta))
		for k, v := range cfg.Meta {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
