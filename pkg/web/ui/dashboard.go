package ui

import (
	g "github.com/maragudk/gomponents"
	htmx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

func DashboardPage() g.Node {
	return Div(
		Div(g.Text("Dashboard")),
		Div(
			UIButton(UIButtonProps{
				Content: g.Text("Logout"),
				Extensions: []g.Node{
					htmx.Post("/auth/logout"),
				},
			}),
		),
	)
}
