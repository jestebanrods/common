package jwt

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewJWEAuthenticator,

	wire.Bind(new(JWEAuthenticator), new(*jweAuthenticator)),
)
