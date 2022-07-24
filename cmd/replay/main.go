package main

import (
	"flag"
	"fmt"
	"log"
	"path"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/replay"
	ws "github.com/lggruspe/polycloze/word_scheduler"
)

type Args struct {
	logFile string
	dbFile  string

	verbose bool
	steps   int // number of reviews to schedule after replay
}

func parseArgs() Args {
	var args Args
	flag.BoolVar(&args.verbose, "v", false, "verbose")
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
	if args.verbose {
		replay.SetVerbosity(true)
	}

	db, err := database.New(args.dbFile)
	if err != nil {
		log.Fatal(err)
	}

	session, err := database.NewSession(db, basedir.Course("eng", inferLanguage(args.logFile)))
	if err != nil {
		log.Fatal(err)
	}

	if err := replay.ReplayFile(session, args.logFile); err != nil {
		log.Fatal(err)
	}

	words, err := ws.GetWordsAt(session, args.steps, replay.Tomorrow())
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
