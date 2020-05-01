package common

import "github.com/sirupsen/logrus"

type LoggerEnv struct {
	Level     string `env:"LOG_LEVEL" envDefault:""`
	Formatter string `env:"LOG_FORMATTER" envDefault:"json"`
}

func NewLogger(env *LoggerEnv) *logrus.Entry {
	logger := logrus.New()

	if env.Formatter == "json" {
		logger.Formatter = &logrus.JSONFormatter{}
	}

	level, err := logrus.ParseLevel(env.Level)
	if err != nil {
		level = logrus.DebugLevel
	}

	logger.Level = level

	return logrus.NewEntry(logger)
}
