package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/flashcards"
)

type Items struct {
	Items []flashcards.Item `json:"items"`
}

func generateFlashcards(db *sql.DB, config Config) func(http.ResponseWriter, *http.Request) {
	ig := flashcards.NewItemGenerator(
		db,
		config.Lang1Db,
		config.Lang2Db,
		config.TranslationDb,
	)
	return func(w http.ResponseWriter, r *http.Request) {
		if config.AllowCORS {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Content-Type", "application/json")

		words, err := ig.GenerateWords(10)
		if err != nil {
			log.Fatal(err)
		}

		items := ig.GenerateItems(words)
		bytes, err := json.Marshal(Items{Items: items})
		if err != nil {
			log.Fatal(err)
		}

		w.Write(bytes)
	}
}

func Mux(config Config) (*http.ServeMux, error) {
	db, err := database.New(config.ReviewDb)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", generateFlashcards(db, config))
	return mux, nil
}
