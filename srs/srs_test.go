package srs

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInitScheduler(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	ws, err := InitWordScheduler(db)

	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}

	if ws.db == nil {
		t.Log("expected WordScheduler.db to be not nil")
		t.Fail()
	}
}
