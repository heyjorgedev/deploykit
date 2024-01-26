package http

import "net/http"

func (s *Server) handlerAuthGetLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.SessionManager.Put(r.Context(), "message", "Hello from a session!")
	}
}
