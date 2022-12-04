package server

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

//go:embed templates
var templates embed.FS

type Template struct {
	templates *template.Template
}

func NewTemplate() *Template {
	t := template.Must(
		template.New("").Funcs(template.FuncMap{}).ParseFS(templates, "templates/*.html"),
	)
	return &Template{templates: t}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
