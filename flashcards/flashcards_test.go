package flashcards

import (
	"testing"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
)

func init() {
	if err := basedir.Init(); err != nil {
		panic(err)
	}
}

func createIg() ItemGenerator {
	db, err := database.New(":memory:")
	if err != nil {
		panic(err)
	}
	return NewItemGenerator(
		db,
		basedir.Language("eng"),
		basedir.Language("spa"),
		basedir.Translation("eng", "spa"),
	)
}

func TestProfiler(t *testing.T) {
	ig := createIg()
	words, err := ig.GenerateWords(10)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	ig.GenerateItems(words)
}
