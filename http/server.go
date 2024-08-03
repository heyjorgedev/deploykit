package http

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/heyjorgedev/deploykit"
)

const ShutdownTimeout = 1 * time.Second

type Server struct {
	ln     net.Listener
	server *http.Server
	router chi.Router

	Addr string

	TeamService deploykit.TeamService
}

func NewServer() *Server {
	r := chi.NewRouter()
	srv := &Server{
		server: &http.Server{},
		router: r,
	}

	// Register global middlewares
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Heartbeat("/up"))

	// Bind the router to the http server
	srv.server.Handler = srv.router

	// Register public directory
	r.Handle("/public/*", handlePublicDir())

	// Register unauthenticated routes
	r.Group(func(r chi.Router) {
		srv.registerInstallWizardRoutes(r)
		srv.registerAuthRoutes(r)
	})

	// Register authenticated routes
	r.Group(func(r chi.Router) {
		srv.registerTeamRoutes(r)
		srv.registerProjectRoutes(r)
	})

	return srv
}

func (s *Server) Open() (err error) {
	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
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
