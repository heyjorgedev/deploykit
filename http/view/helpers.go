package view

import (
	g "github.com/maragudk/gomponents"
	"strings"
)

type nodeTypeDescriber interface {
	Type() g.NodeType
}

// EXPERIMENTAL: MergeAttributes merges all attributes in the given slice of nodes into a single node.
func MergeAttributes(original []g.Node) g.Node {
	var result []g.Node

	for _, n := range original {
		comp, ok := n.(nodeTypeDescriber)
		if !ok {
			continue
		}

		if comp.Type() == g.AttributeType {
			result = append(result, n)
		}

	}

	return g.Group(result)
}

// EXPERIMENTAL: Hacky tailwind merge. Extract to package
func tailwindMerge(origClass string, newClass string) string {
	uniqueMap := map[string]string{}
	attributes := strings.Split(origClass, " ")
	for _, attr := range attributes {
		key := strings.Split(attr, "-")[0]
		uniqueMap[key] = attr
	}

	attributes = strings.Split(newClass, " ")
	for _, attr := range attributes {
		key := strings.Split(attr, "-")[0]
		uniqueMap[key] = attr
	}

	result := []string{}
	for _, attr := range uniqueMap {
		result = append(result, attr)
	}

	return strings.Join(result, " ")
}
