package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/lggruspe/polycloze/buffer"
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
	buf := buffer.NewItemBuffer(ig)
	if err := buf.Fetch(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < n; i++ {
		fmt.Println("i:", i)
		item := buf.Take()
		fmt.Println(item)
	}
}
