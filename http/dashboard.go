package http

import (
	"github.com/heyjorgedev/deploykit/http/view"
	g "github.com/maragudk/gomponents"
	"net/http"
)

func (s *Server) handleDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.LayoutAuth(view.LayoutAuthProps{
			Title:   "Dashboard",
			Content: g.Text("Hello from dashboard!"),
		}).Render(w)
	}
}
