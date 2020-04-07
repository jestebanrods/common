package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"static-server/internal/libraries"

	"static-server/internal/providers"

	"github.com/husobee/vestigo"
	"github.com/sirupsen/logrus"
)

type WebServer interface {
	Shutdown()
	Run(handler http.Handler) error
	RunStaticServer() error
}

type webServer struct {
	halt   chan os.Signal
	logger *logrus.Entry
	http   *http.Server
}

func (s *webServer) shutdownSignal() {
	signal.Notify(s.halt, os.Interrupt, os.Kill, syscall.SIGTERM)
	code := <-s.halt
	s.logger.Infof("interrupt signal detect %s", code)
	s.Shutdown()
}

func (s *webServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.http.Shutdown(ctx)
}

func (s *webServer) Run(handler http.Handler) error {
	go s.shutdownSignal()
	s.http.Handler = handler

	err := s.http.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *webServer) RunStaticServer() error {
	router := vestigo.NewRouter()
	ui := libraries.notFoundWrapper(http.FileServer(http.Dir("/ui")))
	router.Handle("/*", ui)

	s.logger.Info("server on")

	return s.Run(router)
}

func NewWebServer(httpProvider providers.HTTPProvider, loggerProvider providers.LoggerProvider) *webServer {
	server := &webServer{
		http:   httpProvider.HTTPServer(),
		logger: loggerProvider.Logger(),
		halt:   make(chan os.Signal, 1),
	}

	return server
}
