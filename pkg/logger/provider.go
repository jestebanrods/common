package logger

import "github.com/sirupsen/logrus"

type Provider interface {
	Logger() *logrus.Entry
}

type provider struct {
	logger *logrus.Logger
}

func (c *provider) Logger() *logrus.Entry {
	return logrus.NewEntry(c.logger)
}

func NewProvider(env *Env) *provider {
	logger := logrus.New()

	if env.Formatter == "json" {
		logger.Formatter = &logrus.JSONFormatter{}
	}

	level, err := logrus.ParseLevel(env.Level)
	if err == nil {
		logger.Level = level
	} else {
		logger.Level = logrus.DebugLevel
	}

	return &provider{
		logger: logger,
	}
}
