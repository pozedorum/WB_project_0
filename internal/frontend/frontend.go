package frontend

import (
	"html/template"
	"log"
	"net/http"
)

var (
	templateDir = "./internal/frontend/templates"
)

type Frontend struct {
	tmpl *template.Template
}

func New(logger *log.Logger) (*Frontend, error) {
	tmpl, err := template.ParseFiles(templateDir + "/index.html")
	if err != nil {
		return nil, err
	}

	return &Frontend{tmpl: tmpl}, nil
}

func (f *Frontend) RenderIndex(w http.ResponseWriter) error {
	return f.tmpl.ExecuteTemplate(w, "index.html", nil)
}
