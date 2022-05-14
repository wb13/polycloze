package main

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	rs "github.com/lggruspe/polycloze/review_scheduler"
	ws "github.com/lggruspe/polycloze/word_scheduler"
)

func assertNil(value any) {
	if value != nil {
		log.Fatal(value)
	}
}

func main() {
	db, err := ws.New("review.db", "spa.db")
	assertNil(err)

	words, err := ws.GetWords(db, 10)
	assertNil(err)

	for _, word := range words {
		fmt.Println(word)
	}

	if len(words) > 0 {
		assertNil(rs.UpdateReview(db, words[0], true))
	}
}
