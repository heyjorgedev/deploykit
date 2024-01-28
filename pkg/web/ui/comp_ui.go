package ui

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type UIButtonProps struct {
	Type       string
	Content    g.Node
	Href       string
	Extensions []g.Node
	Class      string
}

func UIButton(p UIButtonProps) g.Node {
	classNames := tailwindMerge(
		"flex justify-center rounded-md bg-rose-600 px-5 py-2 text-sm font-semibold text-white shadow-sm hover:bg-rose-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600",
		p.Class,
	)

	if p.Href != "" {
		return A(
			Class(classNames),
			Href(p.Href),
			g.Group(p.Extensions),
			p.Content,
		)
	}

	if p.Type == "" {
		p.Type = "button"
	}

	return Button(
		Class(classNames),
		Type(p.Type),
		g.Group(p.Extensions),
		p.Content,
	)
}
