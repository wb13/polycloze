// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/lggruspe/polycloze/basedir"
	_ "github.com/mattn/go-sqlite3"
)

type Language struct {
	Code          string         `json:"code"` // ISO 639-3
	Name          string         `json:"name"` // in english
	LanguageStats *LanguageStats `json:"stats,omitempty"`
}

// Only used for encoding to json
type Languages struct {
	Languages []Language `json:"languages"`
}

// NOTE Caller has to Close the db.
func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getLanguageName(path string) (string, error) {
	db, err := openDB(path)
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := `select name from language where id = 'l2'`
	row := db.QueryRow(query)

	var name string
	err = row.Scan(&name)
	return name, err
}

// Looks for supported languages in data directories (see basedir package).
func SupportedLanguages() []Language {
	var languages []Language

	targets := make(map[string]string)
	matches, _ := filepath.Glob(path.Join(basedir.DataDir, "[a-z][a-z][a-z]-[a-z][a-z][a-z].db"))
	for _, match := range matches {
		lang := match[len(match)-6 : len(match)-3]
		if _, ok := targets[lang]; ok {
			continue
		}
		name, err := getLanguageName(match)
		if err != nil {
			log.Println(err)
		}
		targets[lang] = name
	}

	for code, name := range targets {
		var lang Language
		lang.Code = code
		lang.Name = name

		if stats, err := getLanguageStats(code); err == nil {
			lang.LanguageStats = stats
		}
		languages = append(languages, lang)
	}
	return languages
}

func languageOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	langs := Languages{Languages: SupportedLanguages()}
	bytes, err := json.Marshal(langs)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := w.Write(bytes); err != nil {
		log.Println(err)
	}
}
