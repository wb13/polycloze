package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/flashcards"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := database.New("review.db")
	if err != nil {
		log.Fatal(err)
	}
	ig := flashcards.NewItemGenerator(
		db,
		"../eng.db",
		"../spa.db",
		"../translations.db",
	)

	items := ig.GenerateItems(10)
	for _, item := range items {
		fmt.Println(item)
	}
}
