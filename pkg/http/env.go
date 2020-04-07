package http

type Env struct {
	Addr               string `env:"HTTP_ADDR" envDefault:"localhost"`
	Port               int    `env:"HTTP_PORT" envDefault:"8080"`
	ReadHeadersTimeout int    `env:"HTTP_READ_HEADERS_TIMEOUT" envDefault:"10"`
	ReadRequestTimeout int    `env:"HTTP_READ_REQUEST_TIMEOUT" envDefault:"20"`
}
