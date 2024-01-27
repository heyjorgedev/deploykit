package http

import (
	"github.com/heyjorgedev/deploykit/http/view"
	g "github.com/maragudk/gomponents"
	"net/http"
	"time"
)

func (s *Server) handleDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Current Date
		view.LayoutAuth(view.LayoutAuthProps{
			Title:   "Dashboard",
			Content: g.Textf("Hello from dashboard!, %d", time.Now().UnixMilli()),
		}).Render(w)
	}
}
