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

func TestInitSchedulerTwice(t *testing.T) {
	// Migration should go smoothly both times, even if there are no changes.
	db, _ := sql.Open("sqlite3", ":memory:")
	if _, err := InitWordScheduler(db); err != nil {
		t.Log("expected err to be nil on first InitWordScheduler", err)
		t.Fail()
	}

	if _, err := InitWordScheduler(db); err != nil {
		t.Log("expected err to be nil on second InitWordScheduler", err)
		t.Fail()
	}
}

// Returns WordScheduler for testing.
func wordScheduler() WordScheduler {
	db, _ := sql.Open("sqlite3", ":memory:")
	ws, _ := InitWordScheduler(db)
	return ws
}

func TestSchedule(t *testing.T) {
	// Result should be empty with no errors.
	ws := wordScheduler()

	words, err := ws.ScheduleNow(100)

	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
	if len(words) > 0 {
		t.Log("expected words to be empty", words)
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	// Only incorrect review needs to be reviewed.
	ws := wordScheduler()

	if err := ws.Update("foo", false); err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
	if err := ws.Update("bar", true); err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}

	words, err := ws.ScheduleNow(100)
	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}

	if len(words) != 1 {
		t.Log("expected different number of results", words)
		t.Fail()
	}
	if words[0] != "foo" {
		t.Log("expected scheduled words to contain \"foo\"", words[0])
		t.Fail()
	}
}

func TestUpdateRepeatedlyCorrect(t *testing.T) {
	// Its streak should increase by one for each update.
	ws := wordScheduler()

	ws.Update("foo", true)
	ws.Update("foo", true)
	ws.Update("foo", true)
	// TODO
}
