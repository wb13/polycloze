// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates/*.html
var templatesFS embed.FS

var templates *template.Template = template.Must(template.ParseFS(templatesFS, "templates/*.html"))

func init() {
	// Check templates.
	names := []string{"home.html", "study.html"}
	for _, name := range names {
		if t := templates.Lookup(name); t == nil {
			log.Fatal("missing template:", name)
		}
	}
}

// Renders template.
// Replies with an internal server error when template execution fails.
// Caller shouldn't make further writes in this case.
func renderTemplate(w http.ResponseWriter, name string, data map[string]any) {
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		log.Println(fmt.Errorf("template execution error: %v", err))
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
}
