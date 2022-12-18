// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

// Writes JSON to file.
func writeJSON(name string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to encode to JSON: %v", err)
	}

	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	if _, err := f.Write(bytes); err != nil {
		return fmt.Errorf("failed to write JSON: %v", err)
	}
	return nil
}

// Parses JSON.
// Writes error to ResponseWriter on error (caller shouldn't write more data).
func parseJSON(w http.ResponseWriter, data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		http.Error(w, "could not parse JSON", http.StatusBadRequest)
		return fmt.Errorf("could not parse JSON: %v", err)
	}
	return nil
}
