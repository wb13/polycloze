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
	"github.com/lggruspe/polycloze/buffer"
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
		path.Join(basedir.DataDir, "translations", "eng-spa.db"),
	)
	buf := buffer.NewItemBuffer(ig, 30)
	for i := 0; i < n; i++ {
		item := buf.Take()
		word := item.Sentence.Parts[1]
		fmt.Println(word, item)
	}

	fmt.Println(":)")
	for i := 0; i < n; i++ {
		item := buf.Take()
		word := item.Sentence.Parts[1]
		fmt.Println(word, item)
	}
}
