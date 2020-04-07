package http

import "github.com/google/wire"

var Set = wire.NewSet(
	NewResponder,
	NewServer,
	wire.Bind(new(Responder), new(*responder)),
	wire.Bind(new(Server), new(*server)),
)
