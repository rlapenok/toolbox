package microservice

import (
	"github.com/rlapenok/toolbox/logger"
	"go.uber.org/zap"
)

type defaultLogerConfig struct{}

func (c *defaultLogerConfig) GetFormat() logger.LogFormat {
	return logger.FormatLogfmt
}

func (c *defaultLogerConfig) GetLevel() string {
	return "debug"
}

func (c *defaultLogerConfig) GetMeta() map[string]any {
	return nil
}

func defaultLogger() *zap.Logger {
	logger, _ := logger.New(&defaultLogerConfig{})
	return logger
}
