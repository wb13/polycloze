package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/flashcards"
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
}

// Only used for encoding to json
type Languages struct {
	Languages []Language `json:"languages"`
}

func getLanguageInfo(path string) (Language, error) {
	var lang Language

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return lang, err
	}
	defer db.Close()

	query := `select code, native, english from info`
	row := db.QueryRow(query)
	if err := row.Scan(&lang.Code, &lang.Native, &lang.English); err != nil {
		return lang, err
	}
	return lang, nil
}

// language: ISO 639-3 code
func languageDatabasePath(language string) string {
	return path.Join(basedir.DataDir, "languages", fmt.Sprintf("%v.db", language))
}

func translationDatabasePath(l1, l2 string) string {
	if l1 == l2 {
		panic("invalid language pair")
	}
	if l2 < l1 {
		l1, l2 = l2, l1
	}
	pair := fmt.Sprintf("%s-%s.db", l1, l2)
	return path.Join(basedir.DataDir, "translations", pair)
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

func loadLanguagePair(l1, l2 string) (*flashcards.ItemGenerator, error) {
	reviewDb := path.Join(basedir.StateDir, "user", fmt.Sprintf("%v.db", l2))
	db, err := database.New(reviewDb)
	if err != nil {
		return nil, err
	}

	ig := flashcards.NewItemGenerator(
		db,
		languageDatabasePath(l1),
		languageDatabasePath(l2),
		translationDatabasePath(l1, l2),
	)
	return &ig, nil
}
