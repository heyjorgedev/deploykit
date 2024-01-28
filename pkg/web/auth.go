package web

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit/pkg/core"
	"github.com/heyjorgedev/deploykit/pkg/web/ui"
	"net/http"
)

type authHandler struct {
	app     core.App
	session *scs.SessionManager
}

func registerAuthRoutes(app core.App, r chi.Router, session *scs.SessionManager) {
	h := &authHandler{
		app:     app,
		session: session,
	}

	r.Get("/auth/login", h.handleGetLogin)
	r.Post("/auth/login", h.handlePostLogin)
}

func (h *authHandler) handleGetLogin(w http.ResponseWriter, r *http.Request) {
	h.session.Put(r.Context(), "user_id", 123)
	_ = ui.LayoutGuest(ui.LayoutGuestProps{
		Title:   "Login",
		Content: ui.AuthLoginForm(ui.AuthLoginFormProps{}),
	}).Render(w)
}

func (h *authHandler) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login"))
}
