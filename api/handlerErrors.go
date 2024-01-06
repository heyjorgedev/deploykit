package api

import "net/http"

func (s *server) handlerErrorNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	}
}
