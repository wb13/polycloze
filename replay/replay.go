// For replaying reviews in log files
package replay

import (
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/logger"
	ws "github.com/lggruspe/polycloze/word_scheduler"
)

func Replay(s *database.Session, events []logger.LogEvent) error {
	for _, event := range events {
		err := ws.UpdateWordAt(s, event.Word, event.Correct, event.Timestamp)
		if err != nil {
			return err
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
