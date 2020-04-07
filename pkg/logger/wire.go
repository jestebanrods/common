package logger

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewLoggerProvider,
	wire.Bind(new(LoggerProvider), new(*loggerProvider)),
)
