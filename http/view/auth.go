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
	return FormEl(Class("grid gap-6"), Method("POST"), Action("/auth/login"), htmx.Post("/auth/login"), htmx.Swap("outerHTML"), htmx.Ext("loading-states"),
		Div(Class("grid gap-2"),
			FormLabel(FormLabelProps{For: "username", Content: g.Text("Username")}),
			FormInput(FormInputProps{
				Name:     "username",
				ID:       "username",
				Value:    p.Username,
				HasError: p.Error != "",
			}),
			g.If(p.Error != "", FormError(g.Text(p.Error))),
		),
		Div(Class("grid gap-2"),
			Label(g.Text("Password"), For("password"), Class("block text-sm font-medium text-gray-700")),
			FormInput(FormInputProps{
				Type: "password",
				Name: "password",
				ID:   "password",
			}),
		),
		Div(
			UIButton(UIButtonProps{
				Type: "submit",
				Content: g.Group([]g.Node{
					g.Text("Login"),
				}),
				Class:      "w-full",
				Extensions: []g.Node{g.Attr("data-loading-disable")},
			}),
		),
	)
}
