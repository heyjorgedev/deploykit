package http

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/heyjorgedev/deploykit"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const ShutdownTimeout = 1 * time.Second

type Server struct {
	ln     net.Listener
	server *http.Server
	// http2Server *http2.Server
	router *chi.Mux

	SessionManager *scs.SessionManager

	Addr string

	ProjectService *deploykit.ProjectService
}

func NewServer() *Server {
	r := chi.NewRouter()
	s := &Server{
		router: r,
	}

	// s.http2Server = &http2.Server{}
	s.server = &http.Server{
		//Handler: h2c.NewHandler(http.HandlerFunc(s.serveHttp), s.http2Server),
		Handler: http.HandlerFunc(s.serveHttp),
	}

	s.registerRoutes()

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
	s.SessionManager.LoadAndSave(s.router).ServeHTTP(w, r)
}
