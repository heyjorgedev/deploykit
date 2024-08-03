package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit/view"
)

func (s *Server) registerProjectRoutes(r chi.Router) {
	r.Get("/projects", s.handleProjectIndex)
}

func (s *Server) handleProjectIndex(w http.ResponseWriter, r *http.Request) {
	view.RenderProjectList(w)
}
