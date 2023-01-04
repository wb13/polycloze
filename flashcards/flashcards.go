// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package flashcards

import (
	"database/sql"
	"fmt"

	"github.com/polycloze/polycloze/database"
	"github.com/polycloze/polycloze/sentences"
	"github.com/polycloze/polycloze/translator"
	"github.com/polycloze/polycloze/word_scheduler"
)

// Different from sentences.Sentence
type Sentence struct {
	ID        int    `json:"id"`    // id in database
	Parts     []Part `json:"parts"` // Odd-numbers parts are blanks
	TatoebaID int64  `json:"tatoebaID,omitempty"`
}

type Item struct {
	Sentence    Sentence               `json:"sentence"`
	Translation translator.Translation `json:"translation"`
}

type ItemGenerator struct {
	db       *sql.DB
	courseDB string // to be attached
}

// Creates an ItemGenerator.
func NewItemGenerator(db *sql.DB, courseDB string) ItemGenerator {
	return ItemGenerator{
		db:       db,
		courseDB: courseDB,
	}
}

func generateItem[T database.Querier](q T, word word_scheduler.Word) (Item, error) {
	var item Item

	sentence, err := sentences.PickSentence(q, word.Word)
	if err != nil {
		return item, err
	}

	translation, err := translator.Translate(q, sentence)
	if err != nil {
		// Panic because this shouldn't happen with generated course files.
		panic(fmt.Errorf("could not translate sentence (%v): %w", sentence, err))
	}
	return Item{
		Translation: translation,
		Sentence: Sentence{
			ID:        sentence.ID,
			Parts:     getParts(sentence.Tokens, word),
			TatoebaID: sentence.TatoebaID,
		},
	}, nil
}

// Creates a cloze item for each word.
func generateItems(con *database.Connection, words []word_scheduler.Word) []Item {
	// To make sure JSON encoding is not nil:
	items := make([]Item, 0)
	for _, word := range words {
		if item, err := generateItem(con, word); err == nil {
			items = append(items, item)
		}
	}
	return items
}

// Returns list of flashcards to show.
// n: max number of flashcards to return.
// Database connection should have access to course and review data.
func Get(
	con *database.Connection,
	n int,
	pred func(word string) bool,
) []Item {
	words, err := word_scheduler.GetWordsWith(con, n, pred)
	if err != nil {
		return nil
	}
	return generateItems(con, words)
}
