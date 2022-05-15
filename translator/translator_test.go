package translator

import (
	"testing"

	"github.com/lggruspe/polycloze/database"
)

var session *database.Session

func init() {
	db, err := database.New(":memory:")
	if err != nil {
		panic(err)
	}
	session, err = database.NewSession(
		db,
		"../cmd/eng.db",
		"../cmd/spa.db",
		"../cmd/translations.db",
	)
	if err != nil {
		panic(err)
	}
}

func TestTranslate(t *testing.T) {
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
