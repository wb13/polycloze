package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lggruspe/polycloze-sentence-picker/sentence_picker"
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

	db, err := sql.Open("sqlite3", ":memory:")
	assertNil(err)

	err = sentence_picker.InitSentencePicker(db, "spa.db", "review.db")
	assertNil(err)

	sentence, err := sentence_picker.PickSentence(db, word)
	assertNil(err)
	fmt.Printf("picked sentence: %v\n", *sentence)
}
