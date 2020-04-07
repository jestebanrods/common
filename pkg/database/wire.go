package database

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewMySQLProvider,

	wire.Bind(new(MySQLProvider), new(*mySQLProvider)),
)
