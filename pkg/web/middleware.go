package web

import (
	"github.com/alexedwards/scs/v2"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func authMiddleware(session *scs.SessionManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := session.GetInt(r.Context(), "userID")
			if userID == 0 {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
