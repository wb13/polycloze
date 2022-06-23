package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/flashcards"
)

func main() {
	if err := basedir.Init(); err != nil {
		log.Fatal(err)
	}

	n := 10
	var err error
	if len(os.Args) >= 2 {
		n, err = strconv.Atoi(os.Args[1])
		if err != nil {
			n = 10
		}
	}

	rand.Seed(time.Now().UnixNano())

	db, err := database.New(basedir.Review("spa"))
	if err != nil {
		log.Fatal(err)
	}
	ig := flashcards.NewItemGenerator(
		db,
		basedir.Language("eng"),
		basedir.Language("spa"),
		basedir.Translation("eng", "spa"),
	)

	start := time.Now()
	words, err := ig.GenerateWords(n)
	if err != nil {
		log.Fatal(err)
	}

	items := ig.GenerateItems(words)
	for _, item := range items {
		fmt.Println(item)
	}

	throughput := float64(len(items)) / time.Since(start).Seconds()
	fmt.Printf("throughput: %v\n", throughput)
}
