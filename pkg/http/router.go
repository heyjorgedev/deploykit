package http

import (
	"github.com/alexedwards/scs/v2"
	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit/pkg/core"
	"github.com/heyjorgedev/deploykit/pkg/http/assets"
	"net/http"
)

func newRouter(app core.App, session *scs.SessionManager) (chi.Router, error) {
	r := chi.NewRouter()
	r.Handle("/assets/*", http.StripPrefix("/assets/", hashfs.FileServer(assets.FS)))

	r.Route("/", func(r chi.Router) {
		r.Use(session.LoadAndSave)
		registerAuthRoutes(app, r)
	})

	return r, nil
}
