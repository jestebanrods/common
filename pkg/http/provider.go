package http

import (
	"fmt"
	"net/http"
	"time"
)

type HTTPEnv struct {
	HTTPAddr               string `env:"HTTP_ADDR" envDefault:"localhost"`
	HTTPPort               int    `env:"HTTP_PORT" envDefault:"9000"`
	HTTPReadHeadersTimeout int    `env:"HTTP_READ_HEADERS_TIMEOUT" envDefault:"10"`
	HTTPReadRequestTimeout int    `env:"HTTP_READ_TIMEOUT" envDefault:"20"`
}

type HTTPProvider interface {
	HTTPAddress() string
	HTTPPort() int
	HTTPAddressPort() string
	HTTPServer() *http.Server
}

type httpProvider struct {
	env *HTTPEnv
}

func (p *httpProvider) HTTPAddress() string {
	return p.env.HTTPAddr
}

func (p *httpProvider) HTTPPort() int {
	return p.env.HTTPPort
}

func (p *httpProvider) HTTPAddressPort() string {
	return fmt.Sprintf("%s:%d", p.env.HTTPAddr, p.env.HTTPPort)
}

func (p *httpProvider) HTTPServer() *http.Server {
	return &http.Server{
		Addr:              p.HTTPAddressPort(),
		ReadHeaderTimeout: time.Duration(p.env.HTTPReadHeadersTimeout) * time.Second,
		ReadTimeout:       time.Duration(p.env.HTTPReadRequestTimeout) * time.Second,
	}
}

func NewHttpProvider(env *HTTPEnv) *httpProvider {
	return &httpProvider{env: env}
}
