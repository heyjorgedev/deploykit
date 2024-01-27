package view

import (
	"github.com/heyjorgedev/deploykit/http/assets"
	g "github.com/maragudk/gomponents"
	htmx "github.com/maragudk/gomponents-htmx"
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
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
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

type LayoutAuthProps struct {
	Title   string
	Content g.Node
}

func LayoutAuth(p LayoutAuthProps) g.Node {
	return Doctype(HTML(Class("h-full bg-gray-50"), Lang("en"),
		Head(
			g.If(p.Title != "", TitleEl(g.Textf("DeployKit - %s", p.Title))),
			g.If(p.Title == "", TitleEl(g.Text("DeployKit"))),
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			Script(Src("https://cdn.tailwindcss.com?plugins=forms")),
			Script(Src(assets.HttpPath("js/htmx@1-9-10.js"))),
			Script(Src(assets.HttpPath("js/htmx@ext-loading-states.js"))),
			Script(Src(assets.HttpPath("js/alpine@3-13-5.js")), Defer()),
		),
		Body(Class("h-full"),
			Div(Class("min-h-full"),
				// NavBar
				Nav(Class("border-b border-gray-200 bg-white"),
					Div(Class("mx-auto max-w-7xl px-4 sm:px-6 lg:px-8"),
						Div(Class("flex h-16 justify-between"),
							Div(Class("flex"),
								Div(Class("flex flex-shrink-0 items-center"),
									Span(Class("text-lg font-bold text-gray-900"), g.Text("DeployKit")),
								),
								Div(Class("hidden sm:-my-px sm:ml-6 sm:flex sm:space-x-8"), htmx.Boost("true"),
									A(Href("/dashboard"), Class("border-transparent text-gray-500 hover:border-rose-600 hover:text-rose-600 inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium"), g.Text("Dashboard")),
									A(Href("/dashboard"), Class("border-transparent text-gray-500 hover:border-rose-600 hover:text-rose-600 inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium"), g.Text("Sites")),
									A(Href("/dashboard"), Class("border-transparent text-gray-500 hover:border-rose-600 hover:text-rose-600 inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium"), g.Text("Databases")),
								),
							),
						),
					),
				),

				// Content
				Div(Class("py-10"),
					Header(
						Div(Class("mx-auto max-w-7xl px-4 sm:px-6 lg:px-8"),
							H1(Class("text-3xl font-bold leading-tight tracking-tight text-gray-900"), g.Text(p.Title)),
						),
					),
					Main(Div(Class("mx-auto max-w-7xl sm:px-6 lg:px-8"), p.Content)),
				),
			),
		),
	))
}
