package http

import (
	"github.com/heyjorgedev/deploykit/http/view"
	g "github.com/maragudk/gomponents"
	"net/http"
)

func (s *Server) handlerAuthGetLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.SessionManager.Put(r.Context(), "message", "Hello from a session!")

		_ = view.LayoutGuest(view.LayoutGuestProps{
			Title:   "Login",
			Content: g.Text("Hello world!"),
		}).Render(w)
	}
}
