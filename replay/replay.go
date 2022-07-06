// For replaying reviews in log files
package replay

import (
	"fmt"
	"time"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/logger"
	ws "github.com/lggruspe/polycloze/word_scheduler"
)

var verbose bool = false
var today time.Time = time.Now().UTC()

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

func Replay(s *database.Session, events []logger.LogEvent) error {
	for _, event := range events {
		err := ws.UpdateWordAt(s, event.Word, event.Correct, event.Timestamp)
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

func ReplayFile(s *database.Session, reviews string) error {
	events, err := logger.ParseFile(reviews)
	if err != nil {
		return err
	}
	return Replay(s, events)
}
