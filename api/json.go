// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// Sends JSON response.
// The caller shouldn't write to w afterwards.
func sendJSON(w http.ResponseWriter, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("failed to encode to JSON:", err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(bytes); err != nil {
		log.Println("failed to send JSON:", err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
}
