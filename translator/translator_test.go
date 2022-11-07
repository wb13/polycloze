// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package translator

import (
	"database/sql"
	"testing"

	"github.com/lggruspe/polycloze/sentences"
	"github.com/lggruspe/polycloze/utils"
)

func translator(reversed bool) *sql.DB {
	db := utils.TestingDatabase()
	populate(db, reversed)
	return db
}

func insertSentence(db *sql.DB, tatoebaID int, text string) {
	query := `
insert into sentence (tatoeba_id, text, tokens, frequency_class)
values (?, ?, '[]', 1)
`
	if _, err := db.Exec(query, tatoebaID, text); err != nil {
		panic(err)
	}
}

func insertTranslation(db *sql.DB, tatoebaID int, text string) {
	query := `insert into translation (tatoeba_id, text) values (?, ?)`
	if _, err := db.Exec(query, tatoebaID, text); err != nil {
		panic(err)
	}
}

func linkSentences(db *sql.DB, source, target int) {
	query := `insert into translates (source, target) values (?, ?)`
	if _, err := db.Exec(query, source, target); err != nil {
		panic(err)
	}
}

// Populates DB with sentences and translations.
func populate(db *sql.DB, reversed bool) {
	if !reversed {
		insertSentence(db, 1, "foo")
		insertTranslation(db, 2, "bar")
		insertTranslation(db, 3, "baz")
		linkSentences(db, 1, 2)
		linkSentences(db, 1, 3)
	} else {
		insertTranslation(db, 1, "foo")
		insertSentence(db, 2, "bar")
		insertSentence(db, 3, "baz")
		linkSentences(db, 2, 1)
		linkSentences(db, 3, 1)
	}
}

func TestTranslate(t *testing.T) {
	t.Parallel()
	// foo -> bar
	session := translator(false)
	defer session.Close()

	sentence, err := sentences.Search(session, "foo")
	if err != nil {
		t.Fatal("sentence not found:", err)
	}

	translation, err := Translate(session, sentence)
	if err != nil {
		t.Fatal("translation failed:", err)
	}
	if len(translation.Text) == 0 {
		t.Fatal("expected translation to be a non-empty string:", translation.Text)
	}
}

func TestReverseTranslate(t *testing.T) {
	t.Parallel()
	// bar -> foo
	session := translator(true)
	defer session.Close()

	sentence, err := sentences.Search(session, "bar")
	if err != nil {
		t.Fatal("sentence not found:", err)
	}

	translation, err := Translate(session, sentence)
	if err != nil {
		t.Fatal("translation failed:", err)
	}
	if len(translation.Text) == 0 {
		t.Fatal("expected translation to be a non-empty string:", translation.Text)
	}
}

func TestTranslateNonTatoebaSentence(t *testing.T) {
	t.Parallel()
	session := translator(false)
	defer session.Close()

	sentence := sentences.Sentence{
		ID:        100,
		TatoebaID: 0, // non-tatoeba sentence <= 0
		Text:      "¿Dónde está la biblioteca?",
		Tokens:    []string{"¿", "Dónde", " ", "está", " ", "la", " ", "biblioteca", "?"},
	}

	_, err := Translate(session, sentence)
	if err == nil {
		t.Fatal("expected translation to fail")
	}
}
