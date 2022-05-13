package main

import (
	"fmt"
	"log"

	"github.com/lggruspe/polycloze-flashcards/flashcards"
)

func main() {
	// TODO where is 'seen' table?
	ig, err := flashcards.NewItemGenerator(
		"review.db",
		"eng.db",
		"spa.db",
		"translations.db",
	)
	if err != nil {
		log.Fatal(err)
	}

	items := ig.GenerateItems(10)
	for _, item := range items {
		fmt.Println(item)
	}
}
