package client

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates/*.html
var templateFS embed.FS

type templates struct {
	t *template.Template
}

func parseTemplates() templates {
	return templates{
		t: template.Must(template.ParseFS(templateFS, "templates/*.html")),
	}
}

func (t templates) execute(w http.ResponseWriter, name string, data any) {
	if err := t.t.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template %q: %v", name, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
