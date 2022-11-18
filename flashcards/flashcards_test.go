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

func BenchmarkGetFlashcards(b *testing.B) {
	db, err := database.New(":memory:")
	if err != nil {
		b.Fatal("expected err to be nil:", err)
	}
	defer db.Close()

	for i := 0; i < b.N; i++ {
		hook := database.AttachCourse(basedir.Course("eng", "deu"))
		Get(db, 10, pred, hook)
	}
}
