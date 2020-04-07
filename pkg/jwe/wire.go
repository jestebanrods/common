package jwe

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewAuthenticator,
	wire.Bind(new(Authenticator), new(*authenticator)),
)
