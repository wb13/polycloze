// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package flashcards

import (
	"testing"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
)

func pred(_ string) bool {
	return true
}

func TestProfiler(t *testing.T) {
	t.Parallel()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	defer db.Close()

	hook := database.AttachCourse(basedir.Course("eng", "spa"))
	Get(db, 10, pred, hook)
}
