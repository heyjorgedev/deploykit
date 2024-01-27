package http

import (
	"fmt"
	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit/http/assets"
	"net/http"
)

func (s *Server) registerRoutes() {
	s.router.Use(s.SessionManager.LoadAndSave)
	s.router.Handle("/assets/*", http.StripPrefix("/assets/", hashfs.FileServer(assets.FS)))

	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		Error(w, http.StatusNotFound)
	})

	// Homepage redirects to auth login if the user is not authenticated
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
	})

	// Authentication routes
	s.router.Route("/auth", func(r chi.Router) {
		r.Get("/login", s.handlerAuthGetLogin())
		r.Post("/login", s.handlerAuthPostLogin())
		r.Route("/mock", func(r chi.Router) {
			r.Use(s.middlewareAuth)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				v := s.SessionManager.GetInt(r.Context(), "userID")
				w.Write([]byte(fmt.Sprintf("Hello from mock! User ID: %d", v)))
			})
		})
	})
}
