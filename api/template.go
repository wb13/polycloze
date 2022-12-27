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

var funcMap template.FuncMap = template.FuncMap{
	"cached": versionedURL,
}

var templates *template.Template = template.Must(
	template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*.html"),
)

// Renders template.
// Replies with an internal server error when template execution fails.
// Caller shouldn't make further writes in this case.
func renderTemplate(w http.ResponseWriter, name string, data map[string]any) {
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		log.Println(fmt.Errorf("template execution error: %v", err))
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
}
