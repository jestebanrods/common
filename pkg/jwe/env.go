package jwe

type Env struct {
	AuthTimeout uint16 `env:"JWE_AUTH_TIMEOUT" envDefault:"60"`
}
