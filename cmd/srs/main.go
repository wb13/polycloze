package main

import (
	"fmt"
	"log"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	ws "github.com/lggruspe/polycloze/word_scheduler"
	_ "github.com/mattn/go-sqlite3"
)

func assertNil(value any) {
	if value != nil {
		log.Fatal(value)
	}
}

func main() {
	db, err := database.New(basedir.Review("spa"))
	assertNil(err)

	session, err := database.NewSession(db, basedir.Course("eng", "spa"))
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
