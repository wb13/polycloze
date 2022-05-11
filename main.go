package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lggruspe/polycloze-sentence-picker/sentence_picker"
)

func assertNil(value any) {
	if value != nil {
		log.Fatal(value)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	assertNil(err)

	err = sentence_picker.InitSentencePicker(db, "spa.db", "review.db")
	assertNil(err)

	sentence, err := sentence_picker.PickSentence(db, "hola")
	assertNil(err)
	fmt.Printf("picked sentence: %v\n", sentence)
}
