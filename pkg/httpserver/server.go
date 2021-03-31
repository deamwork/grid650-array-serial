package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/XSAM/go-hybrid/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mu     sync.Mutex
	router *mux.Router

	addr string
	srv  *http.Server
}

func New(listenAddress string) *HTTPServer {
	return &HTTPServer{
		router: mux.NewRouter().StrictSlash(true),
		addr:   listenAddress,
	}
}

// RegisterRoute append a router to underlying http server mux.
func (s *HTTPServer) RegisterRoute(method string, pattern string, handler http.Handler) *HTTPServer {
	s.mu.Lock()
	s.router.Handle(pattern, handler).Methods(method)
	s.mu.Unlock()
	return s
}

// Addr returns the listening address.
func (s *HTTPServer) Addr() string {
	return s.addr
}

// Serve handles all incoming connections, and,
// the server will blocks util s.Stop() func invoked explicitly.
func (s *HTTPServer) Serve() error {
	if s.srv == nil {
		_ = s.HTTPServer()
	}
	l, err := net.Listen("tcp4", s.addr)
	if err != nil {
		if err != nil {
			log.BgLogger().Fatal("config.rpc.listen",
				zap.String("msg", "Fail to listen on address"),
				zap.Any("addr", s.addr),
				zap.Error(err))
		}
		return errors.Wrap(err, fmt.Sprintf("http server listen failed"))
	}
	s.addr = l.Addr().String()

	return s.srv.Serve(l)
}

// HTTPServer returns the underlying http server object.
func (s *HTTPServer) HTTPServer() *http.Server {
	if s.srv != nil {
		return s.srv
	}
	ns := &http.Server{
		Handler: s.router,
	}
	s.srv = ns
	return s.srv
}

func (s *HTTPServer) GracefulStop() {
	_ = s.Close()
	return
}

func (s *HTTPServer) Close() error {
	return s.srv.Shutdown(context.Background())
}
