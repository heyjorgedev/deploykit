package view

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

type FormInputProps struct {
	Type        string
	Name        string
	ID          string
	Value       string
	HasError    bool
	Placeholder string
}

func FormInput(p FormInputProps) g.Node {
	inputType := p.Type
	if inputType == "" {
		inputType = "text"
	}

	return Input(
		components.Classes{
			"block w-full shadow-sm rounded-md border-0 py-2 ring-1 ring-inset focus:ring-2 focus:ring-inset sm:text-sm sm:leading-6": true,

			"text-red-900 ring-red-300 placeholder:text-red-300 focus:ring-red-500":     p.HasError,
			"text-gray-900 ring-gray-300 placeholder:text-gray-400 focus:ring-rose-600": !p.HasError,
		},
		Type(inputType),
		g.If(p.Name != "", Name(p.Name)),
		g.If(p.ID != "", ID(p.ID)),
		g.If(p.Value != "", Value(p.Value)),
		g.If(p.Placeholder != "", Placeholder(p.Placeholder)),
	)
}

type FormLabelProps struct {
	For     string
	Content g.Node
}

func FormLabel(p FormLabelProps) g.Node {
	return Label(
		Class("block text-sm font-medium leading-6 text-gray-900"),
		For(p.For),
		p.Content,
	)
}

func FormError(children g.Node) g.Node {
	return P(
		Class("text-sm text-red-600"),
		children,
	)
}
