// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func serveFavicons(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	bytes, err := fs.ReadFile(filepath.Join("js", "public", filename))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if _, err := w.Write(bytes); err != nil {
		log.Println(err)
	}
}
