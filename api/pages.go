// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"net/http"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/sessions"
)

func showPage(name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data map[string]any
		db := auth.GetDB(r)
		if s, err := sessions.StartOrResumeSession(db, w, r); err == nil {
			data = s.Data
		}
		if err := renderTemplate(w, name, data); err != nil {
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}
	}
}
