package view

import "io"

func RenderProjectList(w io.Writer) error {
	return htmlTemplates.ExecuteTemplate(w, "project_list.html", nil)
}
