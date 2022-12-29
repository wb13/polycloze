// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/replay"
	ws "github.com/lggruspe/polycloze/word_scheduler"
)

type Args struct {
	logFile string
	dbFile  string

	steps int // number of reviews to schedule after replay
}

func parseArgs() Args {
	var args Args
	flag.IntVar(&args.steps, "n", 0, "number of reviews to schedule after replay")
	flag.Parse()

	nonFlags := flag.Args()
	if len(nonFlags) < 1 {
		log.Fatal("missing arg: path to log file")
	}

	args.logFile = nonFlags[0]
	args.dbFile = ":memory:"

	if len(nonFlags) >= 2 {
		args.dbFile = nonFlags[1]
	}
	return args
}

func main() {
	args := parseArgs()

	db, err := database.OpenReviewDB(args.dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	con, err := database.NewConnection(
		db,
		context.TODO(),
		database.AttachCourse(basedir.Course("eng", inferLanguage(args.logFile))),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := replay.ReplayFile(con, args.logFile); err != nil {
		log.Fatal(err)
	}

	tomorrow := time.Now().Add(24 * time.Hour)
	words, err := ws.GetWordsAt(con, args.steps, tomorrow)
	if err != nil {
		log.Fatal(err)
	}

	if len(words) > 0 {
		fmt.Println("\n# Scheduled for review:")
	}
	for _, word := range words {
		fmt.Println("#", word)
	}
}

func inferLanguage(logFile string) string {
	return path.Base(logFile)[:3]
}
