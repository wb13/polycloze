// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package main

import (
	"context"
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

func pred(_ string) bool {
	return true
}

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

	db, err := database.New(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()

	hook := database.AttachCourse(basedir.Course("eng", "spa"))
	con, err := database.NewConnection(db, context.Background(), hook)
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	items := flashcards.Get(con, n, pred)
	for _, item := range items {
		fmt.Println(item)
	}

	throughput := float64(len(items)) / time.Since(start).Seconds()
	fmt.Printf("throughput: %v\n", throughput)
}
