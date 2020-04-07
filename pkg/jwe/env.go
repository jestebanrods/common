package jwe

import "time"

type Env struct {
	AuthTimeout time.Duration `env:"JWE_AUTH_TIMEOUT" envDefault:"60s"`
}
