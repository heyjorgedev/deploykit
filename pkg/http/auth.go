package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit/http/view"
	"github.com/heyjorgedev/deploykit/pkg/core"
	"net/http"
)

func registerAuthRoutes(app core.App, r chi.Router) {
	h := &authHandler{app: app}

	r.Get("/auth/login", h.handleGetLogin)
	r.Post("/auth/login", h.handlePostLogin)
}

type authHandler struct {
	app core.App
}

func (h *authHandler) handleGetLogin(w http.ResponseWriter, r *http.Request) {
	_ = view.LayoutGuest(view.LayoutGuestProps{
		Title:   "Login",
		Content: view.AuthLoginForm(view.AuthLoginFormProps{}),
	}).Render(w)
}

func (h *authHandler) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login"))
}
