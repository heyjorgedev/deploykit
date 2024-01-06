package api

import "net/http"

func (s *server) handlerAppsList() http.HandlerFunc {
	type app struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Network      string `json:"network"`
		MachineCount int    `json:"machine_count"`
	}
	type response struct {
		TotalApps int   `json:"total_apps"`
		Apps      []app `json:"apps"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		resp := response{
			TotalApps: 0,
			Apps:      make([]app, 0),
		}
		s.respondWithJSON(w, http.StatusOK, resp)
	}
}

func (s *server) handlerAppsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respondWithJSON(w, http.StatusCreated, nil)
	}
}

func (s *server) handlerAppsShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respondWithJSON(w, http.StatusOK, nil)
	}
}

func (s *server) handlerAppsDestroy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respondWithJSON(w, http.StatusOK, nil)
	}
}
