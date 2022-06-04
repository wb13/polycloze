package api

import (
	"database/sql"
	"fmt"
	"path"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/buffer"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/flashcards"
)

type Session struct {
	L1         string
	L2         string
	ItemBuffer *buffer.ItemBuffer
	Database   *sql.DB
}

var globalSession *Session

func changeLanguages(l1 string, l2 string) error {
	if globalSession == nil {
		globalSession = &Session{L1: "", L2: "", ItemBuffer: nil, Database: nil}
	}
	if globalSession.L1 == l1 && globalSession.L2 == l2 {
		return nil
	}

	reviewDb := path.Join(basedir.StateDir, "user", fmt.Sprintf("%v.db", l2))
	db, err := database.New(reviewDb)
	if err != nil {
		return err
	}
	ig := flashcards.NewItemGenerator(
		db,
		languageDatabasePath(l1),
		languageDatabasePath(l2),
		path.Join(basedir.DataDir, "translations.db"),
	)
	buf := buffer.NewItemBuffer(ig, 30)
	globalSession.L1 = l1
	globalSession.L2 = l2
	globalSession.ItemBuffer = &buf

	if globalSession.Database != nil {
		globalSession.Database.Close()
	}

	globalSession.Database = db
	return nil
}
