package main

import (
	"log"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/replay"
)

func main() {
	db, err := database.New("test.db")
	if err != nil {
		log.Fatal(err)
	}

	session, err := database.NewSession(db, ":memory:", ":memory:", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	if err := replay.ReplayFile(session, "test.log"); err != nil {
		log.Fatal(err)
	}
}
