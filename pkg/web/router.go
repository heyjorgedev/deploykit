package web

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/heyjorgedev/deploykit/pkg/core"
	"github.com/heyjorgedev/deploykit/pkg/web/assets"
	"net/http"
)

func newRouter(app core.App, session *scs.SessionManager) (chi.Router, error) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Handle("/assets/*", http.StripPrefix("/assets/", hashfs.FileServer(assets.FS)))

	// Web Routes
	r.Group(func(r chi.Router) {
		// Load session middleware
		r.Use(session.LoadAndSave)

		// TODO: Replace
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			id := session.GetInt(r.Context(), "userID")
			if id == 0 {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}

			w.Write([]byte(fmt.Sprintf("Hello, %d!", id)))
		})

		r.With(authMiddleware(session)).Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("protected"))
		})

		registerAuthRoutes(app, r, session)
	})

	// TODO: API Routes
	r.Route("/api", func(r chi.Router) {})

	return r, nil
}
