package view

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type LayoutGuestProps struct {
	Title   string
	Content g.Node
}

func LayoutGuest(p LayoutGuestProps) g.Node {
	return Doctype(HTML(Class("h-full"), Lang("en"),
		Head(
			g.If(p.Title != "", TitleEl(g.Textf("DeployKit - %s", p.Title))),
			g.If(p.Title == "", TitleEl(g.Text("DeployKit"))),
			Script(Src("https://cdn.tailwindcss.com")),
			Script(
				Src("https://unpkg.com/htmx.org@1.9.10"),
				g.Attr("crossorigin", "anonymous"),
				g.Attr("integrity", "sha384-D1Kt99CQMDuVetoL1lrYwg5t+9QdHe7NLX/SoJYkXDFfX37iInKRy5xLSi8nO7UC"),
			),
			Script(Src("https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"), Defer()),
		),
		Body(
			Class("h-full bg-gray-100"),
			p.Content,
		),
	))
}
