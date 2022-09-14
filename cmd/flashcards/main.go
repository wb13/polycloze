// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

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
	n := 10
	var err error
	if len(os.Args) >= 2 {
		n, err = strconv.Atoi(os.Args[1])
		if err != nil {
			n = 10
		}
	}

	rand.Seed(time.Now().UnixNano())

	db, err := database.New(basedir.Review("eng", "spa"))
	if err != nil {
		log.Fatal(err)
	}
	ig := flashcards.NewItemGenerator(db, basedir.Course("eng", "spa"))

	start := time.Now()

	words, err := ig.GenerateWords(n, func(_ string) bool {
		return true
	})
	if err != nil {
		log.Fatal(err)
	}

	hook := database.AttachCourse(basedir.Course("eng", "spa"))
	items := flashcards.GenerateItems(db, words, hook)
	for _, item := range items {
		fmt.Println(item)
	}

	throughput := float64(len(items)) / time.Since(start).Seconds()
	fmt.Printf("throughput: %v\n", throughput)
}
