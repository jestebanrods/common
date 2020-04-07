package logger

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewProvider,
	wire.Bind(new(Provider), new(*provider)),
)
