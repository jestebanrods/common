package http

import (
	"fmt"
	"net/http"
	"time"
)

type Provider interface {
	Address() string
	Port() int
	AddressPort() string
	Server() *http.Server
}

type provider struct {
	env *Env
}

func (p provider) Address() string {
	return p.env.Addr
}

func (p provider) Port() int {
	return p.env.Port
}

func (p provider) AddressPort() string {
	return fmt.Sprintf("%s:%d", p.env.Addr, p.env.Port)
}

func (p provider) Server() *http.Server {
	return &http.Server{
		Addr:              p.AddressPort(),
		ReadHeaderTimeout: time.Duration(p.env.ReadHeadersTimeout) * time.Second,
		ReadTimeout:       time.Duration(p.env.ReadRequestTimeout) * time.Second,
	}
}

func NewProvider(env *Env) *provider {
	return &provider{env: env}
}
