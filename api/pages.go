// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

//go:embed js/dist/index.* templates/*.html
var fs embed.FS

var templates *template.Template = template.Must(template.ParseFS(fs, "templates/*.html"))

func init() {
	// Check templates.
	names := []string{"home.html", "study.html"}
	for _, name := range names {
		if t := templates.Lookup(name); t == nil {
			log.Fatal("missing template:", name)
		}
	}
}

func serveDist(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	bytes, err := fs.ReadFile(filepath.Join("js", "dist", filename))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	switch filename {
	case "index.js":
		w.Header().Set("Content-Type", "application/javascript")
	case "index.css":
		w.Header().Set("Content-Type", "text/css")
	}

	if _, err := w.Write(bytes); err != nil {
		log.Println(err)
	}
}

func showHome(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "home.html", nil); err != nil {
		log.Println(err)
	}
}
