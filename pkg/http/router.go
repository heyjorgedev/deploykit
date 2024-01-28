package http

import (
	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit/pkg/core"
	"github.com/heyjorgedev/deploykit/pkg/http/assets"
	"net/http"
)

func newRouter(app core.App) (chi.Router, error) {
	r := chi.NewRouter()
	r.Handle("/assets/*", http.StripPrefix("/assets/", hashfs.FileServer(assets.FS)))

	registerAuthRoutes(app, r)

	return r, nil
}
