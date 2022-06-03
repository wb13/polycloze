package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
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

	reviewDb := path.Join(basedir.StateDir, "user", "spa.db")
	db, err := database.New(reviewDb)
	if err != nil {
		log.Fatal(err)
	}
	ig := flashcards.NewItemGenerator(
		db,
		path.Join(basedir.DataDir, "languages", "eng.db"),
		path.Join(basedir.DataDir, "languages", "spa.db"),
		path.Join(basedir.DataDir, "translations.db"),
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
