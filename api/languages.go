package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"path"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lggruspe/polycloze/basedir"
)

func init() {
	if err := basedir.Init(); err != nil {
		log.Fatal(err)
	}
}

type Language struct {
	Code    string `json:"code"` // ISO 639-3
	Native  string `json:"native"`
	English string `json:"english"`

	LanguageStats *LanguageStats `json:"stats,omitempty"`
}

// Only used for encoding to json
type Languages struct {
	Languages []Language `json:"languages"`
}

// NOTE Caller has to Close the db.
func openDb(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getLanguageInfo(path string) (Language, error) {
	var lang Language

	db, err := openDb(path)
	if err != nil {
		return lang, err
	}
	defer db.Close()

	query := `select code, native, english from info`
	row := db.QueryRow(query)
	if err := row.Scan(&lang.Code, &lang.Native, &lang.English); err != nil {
		return lang, err
	}

	if stats, err := getLanguageStats(lang.Code); err == nil {
		lang.LanguageStats = stats
	}
	return lang, nil
}

// Looks for supported languages in data directories (see basedir package).
func SupportedLanguages() []Language {
	var languages []Language
	dir := path.Join(basedir.DataDir, "languages")

	matches, _ := filepath.Glob(path.Join(dir, "[a-z][a-z][a-z].db"))
	for _, match := range matches {
		lang, err := getLanguageInfo(match)
		if err != nil {
			continue
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
	w.Write(bytes)
}
