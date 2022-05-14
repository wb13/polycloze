package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/sentence_picker"
)

func assertNil(value any) {
	if value != nil {
		log.Fatal(value)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing args")
	}
	word := os.Args[1]

	db, err := database.New(":memory:")
	assertNil(err)

	session, err := database.NewSession(db, "../eng.db", "../spa.db", "../translations.db")
	assertNil(err)

	sentence, err := sentence_picker.PickSentence(session, word)
	assertNil(err)
	fmt.Printf("picked sentence: %v\n", *sentence)
}
