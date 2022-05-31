package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/lggruspe/polycloze/buffer"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/flashcards"
	"github.com/lggruspe/polycloze/review_scheduler"
)

type Items struct {
	Items []flashcards.Item `json:"items"`
}

func generateFlashcards(buf *buffer.ItemBuffer, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if 3*len(buf.Channel) <= 2*cap(buf.Channel) {
		go buf.Fetch()
	}

	n := cap(buf.Channel) / 3
	var items []flashcards.Item
	for i := 0; i < n; i++ {
		items = append(items, buf.Take())
	}
	bytes, err := json.Marshal(Items{Items: items})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(bytes)
}

func success() []byte {
	return []byte("{\"success\": true}")
}

func handleReviewUpdate(ig *flashcards.ItemGenerator, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		panic("expected json body in POST request")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic("something went wrong")
	}

	var reviews Reviews
	if err := json.Unmarshal(body, &reviews); err != nil {
		panic("parsing error")
	}

	session, err := ig.Session()
	if err != nil {
		panic("something went wrong")
	}
	defer session.Close()

	for _, review := range reviews.Reviews {
		review_scheduler.UpdateReview(session, review.Word, review.Correct)
	}
	w.Write(success())
}

func createHandler(db *sql.DB, config Config) func(http.ResponseWriter, *http.Request) {
	ig := flashcards.NewItemGenerator(
		db,
		config.Lang1Db,
		config.Lang2Db,
		config.TranslationDb,
	)
	buf := buffer.NewItemBuffer(ig, 30)
	if err := buf.Fetch(); err != nil {
		log.Fatal(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if config.AllowCORS {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
		}

		switch r.Method {
		case "POST":
			handleReviewUpdate(&ig, w, r)
		case "GET":
			generateFlashcards(&buf, w, r)
		}
	}
}

func Mux(config Config) (*http.ServeMux, error) {
	db, err := database.New(config.ReviewDb)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", createHandler(db, config))
	return mux, nil
}
