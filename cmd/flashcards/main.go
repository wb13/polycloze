package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/lggruspe/polycloze/flashcards"
)

func main() {
	rand.Seed(time.Now().UnixNano())

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
