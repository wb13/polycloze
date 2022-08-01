// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package translator

import (
	"testing"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
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
	translation, err := Translate(session, "Hola.")
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
	translation, err := Translate(session, "Hello.")
	if err != nil {
		t.Fatal("translation failed:", err)
	}
	if len(translation.Text) == 0 {
		t.Fatal("expected translation to be a non-empty string:", translation.Text)
	}
}
