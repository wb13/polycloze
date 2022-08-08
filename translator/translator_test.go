// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package translator

import (
	"testing"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/sentences"
)

func newSession(l1, l2 string) *database.Session {
	db, err := database.New(":memory:")
	if err != nil {
		panic(err)
	}

	session, err := database.NewSession(db, basedir.Course(l1, l2))
	if err != nil {
		panic(err)
	}
	return session
}

func TestTranslate(t *testing.T) {
	t.Parallel()
	session := newSession("eng", "spa")

	sentence, err := sentences.Search(session, "Hola.")
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
	session := newSession("spa", "eng")

	sentence, err := sentences.Search(session, "Hello.")
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
	session := newSession("eng", "spa")

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
