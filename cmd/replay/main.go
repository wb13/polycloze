package main

import (
	"flag"
	"log"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/replay"
)

type Args struct {
	logFile string
	dbFile  string

	verbose bool
	steps   int // number of steps to simulate after log file
}

func parseArgs() Args {
	var args Args
	flag.BoolVar(&args.verbose, "v", false, "verbose")
	flag.IntVar(&args.steps, "n", 0, "number of steps to simulate")
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

	session, err := database.NewSession(db, ":memory:", ":memory:", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	if err := replay.ReplayFile(session, args.logFile); err != nil {
		log.Fatal(err)
	}
}
