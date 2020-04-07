package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"assistans-courses/apiserver/internal/providers"

	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
)

type HttpServer interface {
	Shutdown()
	Run(handler http.Handler, middleware ...alice.Constructor) error
	HttpResponder() HttpResponser
}

type httpServer struct {
	halt      chan os.Signal
	logger    *logrus.Entry
	http      *http.Server
	responser *httpResponder
}

func (s *httpServer) shutdownSignal() {
	signal.Notify(s.halt, os.Interrupt, os.Kill, syscall.SIGTERM)
	code := <-s.halt
	s.logger.Infof("interrupt signal detect %s", code)
	s.Shutdown()
}

func (s *httpServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.http.Shutdown(ctx)
}

func (s *httpServer) Run(handler http.Handler, middleware ...alice.Constructor) error {
	go s.shutdownSignal()

	s.http.Handler = alice.New(middleware...).Then(handler)
	s.http.Handler = handler

	err := s.http.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s httpServer) HttpResponder() HttpResponser {
	return s.responser
}

func NewHttpServer(httpProvider providers.HTTPProvider, loggerProvider providers.LoggerProvider) *httpServer {
	return &httpServer{
		http:   httpProvider.HTTPServer(),
		logger: loggerProvider.Logger(),
		halt:   make(chan os.Signal, 1),
	}
}
