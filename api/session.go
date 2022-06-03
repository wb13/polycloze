package api

import (
	"database/sql"
	"path"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/buffer"
	"github.com/lggruspe/polycloze/flashcards"
)

type Session struct {
	L1         string
	L2         string
	ItemBuffer *buffer.ItemBuffer
}

var globalSession *Session

func changeLanguages(db *sql.DB, l1 string, l2 string) error {
	if globalSession == nil {
		globalSession = &Session{L1: "", L2: "", ItemBuffer: nil}
	}
	if globalSession.L1 == l1 && globalSession.L2 == l2 {
		return nil
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
	return nil
}
