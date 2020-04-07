package http

import "github.com/google/wire"

var Set = wire.NewSet(
	NewHttpResponder,
	NewHttpServer,

	wire.Bind(new(HttpResponser), new(*httpResponder)),
	wire.Bind(new(HttpServer), new(*httpServer)),
)
