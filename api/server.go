package api

import "github.com/go-chi/chi/v5"

type server struct {
	router *chi.Mux
}

func NewServer() *server {
	s := &server{
		router: chi.NewRouter(),
	}
	s.routes()
	return s
}
