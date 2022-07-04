package main

import (
	"log"
	"os"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/replay"
)

func parseArgs() (string, string) {
	if len(os.Args) < 2 {
		log.Fatal("missing arg: path to log file")
	}

	logFile := os.Args[1]
	dbFile := ":memory:"

	if len(os.Args) >= 3 {
		dbFile = os.Args[2]
	}
	return logFile, dbFile
}

func main() {
	logFile, dbFile := parseArgs()

	db, err := database.New(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	session, err := database.NewSession(db, ":memory:", ":memory:", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	if err := replay.ReplayFile(session, logFile); err != nil {
		log.Fatal(err)
	}
}
