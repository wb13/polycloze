// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package flashcards

import (
	"testing"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/word_scheduler"
)

func createIg() ItemGenerator {
	db, err := database.New(":memory:")
	if err != nil {
		panic(err)
	}
	return NewItemGenerator(db, basedir.Course("eng", "spa"))
}

func TestProfiler(t *testing.T) {
	t.Parallel()

	ig := createIg()

	session, err := ig.Session()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	defer session.Close()

	words, err := word_scheduler.GetWordsWith(session, 10, func(_ string) bool {
		return true
	})
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	ig.GenerateItems(words)
}
