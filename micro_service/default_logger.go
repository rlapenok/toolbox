package microservice

import (
	"github.com/rlapenok/toolbox/logger"
)

func defaultLogger() logger.Config {
	config := logger.Config{
		ConsoleEnabled: true,
		ConsoleFormat:  logger.FormatLogfmt,
		ConsoleMode:    logger.ConsoleSplit,
		ConsoleLevel:   "debug",

		FileEnabled:    true,
		FilePath:       "logs/app.log",
		FileFormat:     logger.FormatLogfmt,
		FileLevel:      "debug",
		FileMaxSizeMB:  100,
		FileMaxBackups: 10,
		FileMaxAgeDays: 30,
		FileCompress:   true,
	}

	return config
}
