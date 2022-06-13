package translator

import (
	"fmt"
	"path"
	"testing"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
)

func newSession(l1, l2 string) *database.Session {
	if err := basedir.Init(); err != nil {
		panic(err)
	}

	db, err := database.New(":memory:")
	if err != nil {
		panic(err)
	}

	pair1, pair2 := l1, l2
	if l2 < l1 {
		pair1, pair2 = l2, l1
	}
	session, err := database.NewSession(
		db,
		path.Join(basedir.DataDir, "languages", fmt.Sprintf("%s.db", l1)),
		path.Join(basedir.DataDir, "languages", fmt.Sprintf("%s.db", l2)),
		path.Join(basedir.DataDir, "translations", fmt.Sprintf("%s-%s.db", pair1, pair2)),
	)
	if err != nil {
		panic(err)
	}
	return session
}

func TestTranslate(t *testing.T) {
	session := newSession("eng", "spa")
	translation, err := Translate(session, "Hola.")
	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
	if len(translation) == 0 {
		t.Log("expected translation to be a non-empty string", translation)
		t.Fail()
	}
}

func TestReverseTranslate(t *testing.T) {
	session := newSession("spa", "eng")
	translation, err := Translate(session, "Hello.")
	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
	if len(translation) == 0 {
		t.Log("expected translation to be a non-empty string", translation)
		t.Fail()
	}
}
