package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit"
)

const ShutdownTimeout = 1 * time.Second

type Server struct {
	ln     net.Listener
	server *http.Server
	// http2Server *http2.Server
	router *chi.Mux

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

	r.NotFound(s.handleNotFound())

	return s
}

func (s *Server) handlerFunc(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err != nil {
			switch err.(type) {
			case *deploykit.Error:
				s.Error(w, r, deploykit.ErrorMessage(err), http.StatusBadRequest)
			default:
				s.Error(w, r, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func (s *Server) Error(w http.ResponseWriter, r *http.Request, message string, status int) {
	s.respond(w, r, status, Envelope[any]{
		Success: false,
		Message: message,
	})
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

func (s *Server) respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) error {
	b := bytes.NewBuffer(nil)

	if err := json.NewEncoder(b).Encode(data); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b.Bytes())
	return nil
}

func (s *Server) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Error(w, r, "not found", http.StatusNotFound)
	}
}
