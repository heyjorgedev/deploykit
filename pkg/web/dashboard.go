package web

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/heyjorgedev/deploykit/pkg/core"
	"github.com/heyjorgedev/deploykit/pkg/web/ui"
	"net/http"
)

type dashboardHandler struct {
	app     core.App
	session *scs.SessionManager
}

func registerDashboardRoutes(app core.App, r chi.Router, session *scs.SessionManager) {
	h := &dashboardHandler{app: app, session: session}

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware(session))
		r.Get("/dashboard", h.handleGetDashboard)
	})
}

func (h *dashboardHandler) handleGetDashboard(w http.ResponseWriter, r *http.Request) {
	_ = ui.LayoutAuth(ui.LayoutAuthProps{
		Title:   "Dashboard",
		Content: ui.DashboardPage(),
	}).Render(w)
}
