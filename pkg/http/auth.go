package http

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit/http/view"
	"github.com/heyjorgedev/deploykit/pkg/core"
	"net/http"
)

func registerAuthRoutes(app core.App, r chi.Router, session *scs.SessionManager) {
	h := &authHandler{app: app, session: session}

	r.Get("/auth/login", h.handleGetLogin)
	r.Post("/auth/login", h.handlePostLogin)
}

type authHandler struct {
	app     core.App
	session *scs.SessionManager
}

func (h *authHandler) handleGetLogin(w http.ResponseWriter, r *http.Request) {
	h.session.Put(r.Context(), "user_id", 123)
	_ = view.LayoutGuest(view.LayoutGuestProps{
		Title:   "Login",
		Content: view.AuthLoginForm(view.AuthLoginFormProps{}),
	}).Render(w)
}

func (h *authHandler) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login"))
}
