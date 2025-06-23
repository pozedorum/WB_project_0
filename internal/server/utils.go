package server

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gomarkdown/markdown"
)

func renderTemplate(w http.ResponseWriter, templatePath string, data interface{}) {
	tmpl, err := template.New(filepath.Base(templatePath)).
		Funcs(template.FuncMap{
			"markdown": renderMarkdown,
		}).ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Template not found: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template render error: "+err.Error(), http.StatusInternalServerError)
	}
}

func renderMarkdown(text string) template.HTML {
	return template.HTML(markdown.ToHTML([]byte(text), nil, nil))
}
