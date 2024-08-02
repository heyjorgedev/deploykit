package view

import (
	"embed"
	"html/template"
	"io"
)

//go:embed *.html
var viewFS embed.FS
var htmlTemplates *template.Template

func init() {
	htmlTemplates = template.Must(template.ParseFS(viewFS, "*.html"))
}

func RenderIndex(w io.Writer) error {
	htmlTemplates.ExecuteTemplate(w, "login.html", nil)
	return nil
}

func RenderLoginForm(w io.Writer) error {
	htmlTemplates.ExecuteTemplate(w, "login_form", nil)
	return nil
}
