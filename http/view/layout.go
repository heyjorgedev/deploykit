package view

import (
	"github.com/heyjorgedev/deploykit/http/assets"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type LayoutGuestProps struct {
	Title   string
	Content g.Node
}

func LayoutGuest(p LayoutGuestProps) g.Node {
	return Doctype(HTML(Class("h-full bg-gray-50"), Lang("en"),
		Head(
			g.If(p.Title != "", TitleEl(g.Textf("DeployKit - %s", p.Title))),
			g.If(p.Title == "", TitleEl(g.Text("DeployKit"))),
			Script(Src("https://cdn.tailwindcss.com?plugins=forms")),
			Script(Src(assets.HttpPath("js/htmx@1-9-10.js"))),
			Script(Src(assets.HttpPath("js/htmx@ext-loading-states.js"))),
			Script(Src(assets.HttpPath("js/alpine@3-13-5.js")), Defer()),
		),
		Body(Class("h-full antialised bg-gray-50 text-gray-950 flex items-center flex-grow"),
			Div(Class("w-full sm:max-w-lg mx-auto bg-white px-12 py-12 shadow-sm ring-1 ring-gray-950/5 sm:rounded-xl"),
				Div(Class("flex-col items-center pb-10"),
					H1(Class("text-lg font-bold text-center text-gray-950 mb-4"), g.Text("DeployKit")),
					H2(Class("text-2xl font-bold text-center text-gray-950"), g.Text("Sign In")),
				),
				p.Content,
			),
		),
	))
}
