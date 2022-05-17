package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/flashcards"
)

func main() {
	n := 10
	var err error
	if len(os.Args) >= 2 {
		n, err = strconv.Atoi(os.Args[1])
		if err != nil {
			n = 10
		}
	}

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

	start := time.Now()

	items := ig.GenerateItems(n)
	for _, item := range items {
		fmt.Println(item)
	}

	throughput := float64(len(items)) / time.Since(start).Seconds()
	fmt.Printf("throughput: %v\n", throughput)
}
