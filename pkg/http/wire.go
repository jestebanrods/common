package http

import "github.com/google/wire"

var Set = wire.NewSet(
	NewResponder,
	NewServer,
	NewProvider,
	wire.Bind(new(Responder), new(*responder)),
	wire.Bind(new(Server), new(*server)),
	wire.Bind(new(Provider), new(*provider)),
)
