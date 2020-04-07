package logger

import "github.com/sirupsen/logrus"

type LoggerEnv struct {
	LogLevel     string `env:"LOG_LEVEL" envDefault:""`
	LogFormatter string `env:"LOG_FORMATTER" envDefault:"json"`
}

type LoggerProvider interface {
	Logger() *logrus.Entry
}

type loggerProvider struct {
	logger *logrus.Logger
}

func (c *loggerProvider) Logger() *logrus.Entry {
	return logrus.NewEntry(c.logger)
}

func NewLoggerProvider(env *LoggerEnv) *loggerProvider {
	logger := logrus.New()

	if env.LogFormatter == "json" {
		logger.Formatter = &logrus.JSONFormatter{}
	}

	level, err := logrus.ParseLevel(env.LogLevel)
	if err == nil {
		logger.Level = level
	} else {
		logger.Level = logrus.DebugLevel
	}

	return &loggerProvider{
		logger: logger,
	}
}
