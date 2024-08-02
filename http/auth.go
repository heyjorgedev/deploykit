package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit/view"
)

func (s *Server) registerAuthRoutes(r chi.Router) {

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		view.RenderIndex(w)
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		hxRequest := r.Header.Get("HX-Request")

		if hxRequest == "true" {
			view.RenderLoginForm(w)
			return
		}

		s.redirectTemporary(w, r, "/login")
	})
}
