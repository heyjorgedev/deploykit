package view

import (
	g "github.com/maragudk/gomponents"
	htmx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

type AuthLoginFormProps struct {
	Username string
	Error    string
}

func AuthLoginForm(p AuthLoginFormProps) g.Node {
	return FormEl(Class("grid gap-6"), Method("POST"), Action("/auth/login"), htmx.Post("/auth/login"), htmx.Swap("outerHTML"),
		g.If(p.Error != "", Div(Class("bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative"), g.Attr("role", "alert"), g.Text(p.Error))),
		Div(
			Label(g.Text("Username"), For("username"), Class("block text-sm font-medium text-gray-700")),
			Input(Class("mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50"), Type("text"), Name("username"), ID("username"), Value(p.Username)),
		),
		Div(
			Label(g.Text("Password"), For("password"), Class("block text-sm font-medium text-gray-700")),
			Input(Class("mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50"), Type("password"), Name("password"), ID("password")),
		),
		Div(
			UIButton(UIButtonProps{
				Type:    "submit",
				Content: g.Text("Login"),
				Class:   "w-full",
			}),
		),
	)
}
