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
	"github.com/jorgemurta/deploykit"
)

const ShutdownTimeout = 1 * time.Second

type Server struct {
	ln     net.Listener
	server *http.Server
	// http2Server *http2.Server
	router *chi.Mux

	Addr string

	AppService deploykit.AppService
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

	r.Route("/apps", func(r chi.Router) {
		r.Get("/", s.handleAppsList().ServeHTTP)
		r.Post("/", s.handleAppsStore().ServeHTTP)
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

func (s *Server) responseWithValidationErrors(w http.ResponseWriter, r *http.Request, errs deploykit.ValidationErrors) error {
	return s.respond(w, r, http.StatusBadRequest, ResourceResponse[any]{
		Errors: errs,
	})
}

func (s *Server) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not found"))
	}
}

func (s *Server) handleAppsList() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		apps, err := s.AppService.FindAll(r.Context())
		if err != nil {
			return err
		}

		return s.respond(w, r, http.StatusOK, ResourceResponse[[]*deploykit.App]{Data: apps})
	}
}

func (s *Server) handleAppsStore() HandlerFunc {
	type Request struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}

		app := &deploykit.App{
			Name: req.Name,
		}

		err := app.Validate()
		if err != nil {
			return s.responseWithValidationErrors(w, r, *err)
		}

		if err := s.AppService.Create(r.Context(), app); err != nil {
			return err
		}

		return s.respond(w, r, http.StatusCreated, ResourceResponse[*deploykit.App]{Data: app})
	}
}
