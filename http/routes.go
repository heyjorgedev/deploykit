package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit/http/view"
	"net/http"
)

func (s *Server) registerRoutes() {
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		_ = view.NotFoundPage().Render(w)
	})

	// Homepage redirects to auth login if the user is not authenticated
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
	})

	// Authentication routes
	s.router.Route("/auth", func(r chi.Router) {
		r.Get("/login", s.handlerAuthGetLogin())
	})
}
