package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestUpgrade(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	if err := Upgrade(db); err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
}

func TestUpgradeTwice(t *testing.T) {
	// Migration should go smoothly both times, even if there are no changes.
	db, _ := sql.Open("sqlite3", ":memory:")
	if err := Upgrade(db); err != nil {
		t.Log("expected err to be nil on first upgrade", err)
		t.Fail()
	}
	if err := Upgrade(db); err != nil {
		t.Log("expected err to be nil on second upgrade", err)
		t.Fail()
	}
}
