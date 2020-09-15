package hecho

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

// Template is very simple template engine
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
