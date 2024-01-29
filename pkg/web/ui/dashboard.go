package ui

import (
	g "github.com/maragudk/gomponents"
	htmx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

type DashboardPageProps struct {
	Name string
}

func DashboardPage(p DashboardPageProps) g.Node {
	return Div(
		Div(g.Text("Dashboard")),
		Div(g.Textf("Hello, %s!", p.Name)),
		Div(
			UIButton(UIButtonProps{
				Content: g.Text("Logout"),
				Extensions: []g.Node{
					htmx.Post("/auth/logout"),
				},
			}),
			UIButton(UIButtonProps{
				Content: g.Text("Go to Google"),
				Href:    "https://google.com",
				Class:   "justify-end",
			}),
		),
	)
}
