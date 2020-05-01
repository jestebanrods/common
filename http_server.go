package common

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/husobee/vestigo"
	"github.com/justinas/alice"
)

type HTTPEnv struct {
	Addr               string `env:"HTTP_ADDR" envDefault:"localhost"`
	Port               int    `env:"HTTP_PORT" envDefault:"8080"`
	ReadHeadersTimeout int    `env:"HTTP_READ_HEADERS_TIMEOUT" envDefault:"10"`
	ReadRequestTimeout int    `env:"HTTP_READ_REQUEST_TIMEOUT" envDefault:"20"`
}

type Server struct {
	halt chan os.Signal
	http *http.Server
}

func (s *Server) Run(handler http.Handler, middleware ...alice.Constructor) error {
	go s.shutdownSignal()

	s.http.Handler = alice.New(middleware...).Then(handler)
	s.http.Handler = handler

	err := s.http.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.http.Shutdown(ctx)
}

func (s *Server) RunStaticServer(dir string) error {
	router := vestigo.NewRouter()
	router.Handle("/*", FrontendHandler(dir))
	return s.Run(router)
}

func (s *Server) shutdownSignal() {
	signal.Notify(s.halt, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-s.halt
	s.Shutdown()
}

func FrontendHandler(publicDir string) http.HandlerFunc {
	handler := http.FileServer(http.Dir(publicDir))

	return func(w http.ResponseWriter, req *http.Request) {
		urlPath := req.URL.Path
		if strings.Contains(urlPath, ".") || urlPath == "/" {
			handler.ServeHTTP(w, req)
			return
		}

		http.ServeFile(w, req, path.Join(publicDir, "/index.html"))
	}
}

func NewServer(env *HTTPEnv) *Server {
	return &Server{
		halt: make(chan os.Signal, 1),
		http: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", env.Addr, env.Port),
			ReadHeaderTimeout: time.Duration(env.ReadHeadersTimeout) * time.Second,
			ReadTimeout:       time.Duration(env.ReadRequestTimeout) * time.Second,
		},
	}
}
