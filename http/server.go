package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jorgemurta/deploykit"
)

const ShutdownTimeout = 1 * time.Second

type Server struct {
	ln     net.Listener
	server *http.Server
	router *chi.Mux

	Addr string

	AppService deploykit.AppService
}

func NewServer() *Server {
	r := chi.NewRouter()
	s := &Server{
		server: &http.Server{},
		router: r,
	}
	s.server.Handler = http.HandlerFunc(s.serveHttp)

	r.NotFound(s.handleNotFound())

	return s
}

func (s *Server) Open() (err error) {
	if s.Addr == "" {
		return fmt.Errorf("addr required")
	}

	s.ln, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	go s.server.Serve(s.ln)

	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) serveHttp(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not found"))
	}
}
