package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/husobee/vestigo"
	"github.com/justinas/alice"
)

type Server interface {
	Shutdown()
	Run(handler http.Handler, middleware ...alice.Constructor) error
	RunStaticServer(dir string) error
}

type server struct {
	halt chan os.Signal
	http *http.Server
}

func (s *server) shutdownSignal() {
	signal.Notify(s.halt, os.Interrupt, os.Kill, syscall.SIGTERM)
	s.Shutdown()
}

func (s *server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.http.Shutdown(ctx)
}

func (s *server) Run(handler http.Handler, middleware ...alice.Constructor) error {
	go s.shutdownSignal()

	s.http.Handler = alice.New(middleware...).Then(handler)
	s.http.Handler = handler

	err := s.http.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *server) RunStaticServer(dir string) error {
	router := vestigo.NewRouter()
	ui := NotFoundWrapper(http.FileServer(http.Dir(dir)))
	router.Handle("/*", ui)
	return s.Run(router)
}

func NewServer(httpProvider Provider) *server {
	return &server{
		http: httpProvider.Server(),
		halt: make(chan os.Signal, 1),
	}
}
