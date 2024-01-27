package http

import (
	"github.com/heyjorgedev/deploykit/http/view"
	"net/http"
)

func (s *Server) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.SessionManager.Exists(r.Context(), "userID") {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) handlerAuthGetLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = view.LayoutGuest(view.LayoutGuestProps{
			Title:   "Login",
			Content: view.AuthLoginForm(view.AuthLoginFormProps{}),
		}).Render(w)
	}
}

func (s *Server) handlerAuthPostLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			Error(w, http.StatusBadRequest)
			return
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")
		user, err := s.AuthService.AttemptCredentials(username, password)
		if err != nil {
			view.AuthLoginForm(view.AuthLoginFormProps{
				Error:    "The credentials you provided are incorrect.",
				Username: username,
			}).Render(w)
			return
		}

		err = s.SessionManager.RenewToken(r.Context())
		if err != nil {
			Error(w, http.StatusInternalServerError)
			return
		}

		s.SessionManager.Put(r.Context(), "userID", user.ID)
		Redirect(w, r, "/auth/mock", http.StatusSeeOther)
	}
}
