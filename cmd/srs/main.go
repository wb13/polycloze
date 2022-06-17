package main

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lggruspe/polycloze/database"
	ws "github.com/lggruspe/polycloze/word_scheduler"
)

func assertNil(value any) {
	if value != nil {
		log.Fatal(value)
	}
}

func main() {
	db, err := database.New("review.db")
	assertNil(err)

	// TODO
	session, err := database.NewSession(db, "eng.db", "spa.db", "")
	assertNil(err)

	words, err := ws.GetWords(session, 10)
	assertNil(err)

	for _, word := range words {
		fmt.Println(word)
	}

	if len(words) > 0 {
		assertNil(ws.UpdateWord(session, words[0], true))
	}
}
