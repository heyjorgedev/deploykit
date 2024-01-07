package http

import (
	"context"
	"encoding/json"
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

	r.Route("/apps", func(r chi.Router) {
		r.Get("/", s.handleAppsList())
	})

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

func (s *Server) respond(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		// TODO: Maybe change to a buffer to avoid partial responses.
		json.NewEncoder(w).Encode(data)
	}
}

func (s *Server) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not found"))
	}
}

func (s *Server) handleAppsList() http.HandlerFunc {
	type Response struct {
		Data []*deploykit.App `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		apps, err := s.AppService.FindAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.respond(w, http.StatusOK, Response{Data: apps})
	}
}
