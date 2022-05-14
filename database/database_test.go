//go:build sqlite_math_functions

package database

import (
	"database/sql"
	"path"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestUpgrade(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	err := Upgrade(db, path.Join("migrations", "review_scheduler"))
	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
}

func TestUpgradeTwice(t *testing.T) {
	// Migration should go smoothly both times, even if there are no changes.
	db, _ := sql.Open("sqlite3", ":memory:")
	migrations := path.Join("migrations", "review_scheduler")
	if err := Upgrade(db, migrations); err != nil {
		t.Log("expected err to be nil on first upgrade", err)
		t.Fail()
	}
	if err := Upgrade(db, migrations); err != nil {
		t.Log("expected err to be nil on second upgrade", err)
		t.Fail()
	}
}
