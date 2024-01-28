package web

import (
	"fmt"
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
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.app.Dao().FindUserByUsername(username)
	if err != nil {
		ui.AuthLoginForm(ui.AuthLoginFormProps{
			Username: username,
			Error:    "Credentials are wrong",
		}).Render(w)
		return
	}

	if isValidPassword := user.ValidatePassword(password); !isValidPassword {
		ui.AuthLoginForm(ui.AuthLoginFormProps{
			Username: username,
			Error:    "Credentials are wrong",
		}).Render(w)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s", user.Name)))
}
