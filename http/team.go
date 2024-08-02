package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit"
	"github.com/heyjorgedev/deploykit/http/htmx"
	"github.com/heyjorgedev/deploykit/view"
)

func (s *Server) registerTeamRoutes(r chi.Router) {
	r.Get("/teams", s.handleTeamIndex)
	r.Post("/teams", s.handleTeamCreate)
}

func (s *Server) handleTeamIndex(w http.ResponseWriter, r *http.Request) {
	teams, _, err := s.TeamService.FindTeams(r.Context(), deploykit.TeamFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var selectedTeam *deploykit.Team
	if len(teams) > 0 {
		selectedTeam = teams[0]
	}

	view.RenderTeamList(w, view.TeamListData{
		SelectedTeam: selectedTeam,
		Teams:        teams,
	})
}

func (s *Server) handleTeamCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")

	team := &deploykit.Team{
		Name: name,
	}

	// Save team to database
	err := s.TeamService.CreateTeam(r.Context(), team)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if htmx.IsHTMXRequest(w, r) {
		w.Header().Set("Content-Type", "text/html")
		htmx.Trigger(w, r, htmx.TriggerPayload{
			"team-created": "Team created",
		})
		w.WriteHeader(201)
		view.RenderTeamListItem(w, team)
		return
	}

	http.Redirect(w, r, "/teams", http.StatusSeeOther)
}
