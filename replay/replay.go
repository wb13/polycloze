// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// For replaying reviews in log files
package replay

import (
	"fmt"
	"time"

	"github.com/lggruspe/polycloze/database"
	ws "github.com/lggruspe/polycloze/word_scheduler"
)

var (
	verbose bool      = false
	today   time.Time = time.Now().UTC()
)

func SetVerbosity(verbosity bool) {
	verbose = verbosity
}

func Today() time.Time {
	return today.UTC()
}

func Tomorrow() time.Time {
	day, _ := time.ParseDuration("24h")
	return Today().Add(day).UTC()
}

func Replay(c *database.Connection, events []LogEvent) error {
	for _, event := range events {
		err := ws.UpdateWordAt(c, event.Word, event.Correct, event.Timestamp)
		if err != nil {
			return err
		}

		today = event.Timestamp

		if verbose {
			fmt.Println(event)
		}
	}
	return nil
}

func ReplayFile(c *database.Connection, reviews string) error {
	events, err := ParseFile(reviews)
	if err != nil {
		return err
	}
	return Replay(c, events)
}
