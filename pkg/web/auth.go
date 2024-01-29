package web

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/heyjorgedev/deploykit/pkg/core"
	"github.com/heyjorgedev/deploykit/pkg/web/ui"
	"net/http"
	"time"
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

	// Guest Routes
	r.Group(func(r chi.Router) {
		r.Use(guestMiddleware(session))

		// Login
		r.Get("/auth/login", h.handleGetLogin)
		r.With(httprate.LimitByIP(30, 1*time.Minute)).Post("/auth/login", h.handlePostLogin)

	})

	// Authenticated Routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware(session))
		r.Post("/auth/logout", h.handlePostLogout)
	})
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

	if err := h.session.RenewToken(r.Context()); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	h.session.Put(r.Context(), "userID", user.ID)

	if isHTMXRequest(r) {
		htmxRedirect(w, r, "/", http.StatusCreated)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *authHandler) handlePostLogout(w http.ResponseWriter, r *http.Request) {
	if err := h.session.RenewToken(r.Context()); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	h.session.Remove(r.Context(), "userID")

	if isHTMXRequest(r) {
		htmxRedirect(w, r, "/", http.StatusAccepted)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
