package api

import "github.com/go-chi/chi/v5"

func (s *server) routes() {
	s.router.NotFound(s.handlerErrorNotFound())

	s.router.Route("/apps", func(r chi.Router) {
		r.Get("/", s.handlerAppsList())
		r.Post("/", s.handlerAppsStore())
		r.Get("/{id}", s.handlerAppsShow())
		r.Delete("/{id}", s.handlerAppsDestroy())
	})
}
