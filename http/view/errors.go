package view

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
	"net/http"
)

func ErrorPage(statusCode int) g.Node {
	return c.HTML5(c.HTML5Props{
		Title: "DeployKit - Not Found",
		Head: []g.Node{
			Script(Src("https://cdn.tailwindcss.com")),
		},
		Body: []g.Node{
			Class("bg-gray-100 min-h-screen flex items-center justify-center text-lg text-gray-500 uppercase tracking-wider"),
			g.Textf("%d | %s", statusCode, http.StatusText(statusCode)),
		},
	})
}

func NotFoundPage() g.Node {
	return ErrorPage(http.StatusNotFound)
}
