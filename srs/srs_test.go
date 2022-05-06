//go:build sqlite_math_functions

package srs

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInitScheduler(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	ws, err := InitReviewScheduler(db)

	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}

	if ws.db == nil {
		t.Log("expected ReviewScheduler.db to be not nil")
		t.Fail()
	}
}

func TestInitSchedulerTwice(t *testing.T) {
	// Migration should go smoothly both times, even if there are no changes.
	db, _ := sql.Open("sqlite3", ":memory:")
	if _, err := InitReviewScheduler(db); err != nil {
		t.Log("expected err to be nil on first InitReviewScheduler", err)
		t.Fail()
	}

	if _, err := InitReviewScheduler(db); err != nil {
		t.Log("expected err to be nil on second InitReviewScheduler", err)
		t.Fail()
	}
}

// Returns ReviewScheduler for testing.
func reviewScheduler() ReviewScheduler {
	db, _ := sql.Open("sqlite3", ":memory:")
	ws, _ := InitReviewScheduler(db)
	return ws
}

func TestSchedule(t *testing.T) {
	// Result should be empty with no errors.
	ws := reviewScheduler()

	items, err := ws.ScheduleNow(100)

	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
	if len(items) > 0 {
		t.Log("expected items to be empty", items)
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	// Only incorrect review needs to be reviewed.
	ws := reviewScheduler()

	if err := ws.Update("foo", false); err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
	if err := ws.Update("bar", true); err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}

	items, err := ws.ScheduleNow(100)
	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}

	if len(items) != 1 {
		t.Log("expected different number of results", items)
		t.Fail()
	}
	if items[0] != "foo" {
		t.Log("expected scheduled items to contain \"foo\"", items[0])
		t.Fail()
	}
}

func TestUpdateRecentlyAnsweredItemDoesntGetScheduled(t *testing.T) {
	ws := reviewScheduler()
	items := []string{"foo", "bar", "baz"}
	for _, item := range items {
		ws.Update(item, true)
	}

	items, err := ws.ScheduleNow(-1)
	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
	if len(items) > 0 {
		t.Log("expected items to be empty", items)
		t.Fail()
	}
}

func TestUpdateRepeatedlyCorrect(t *testing.T) {
	ws := reviewScheduler()
	ws.Update("foo", true)
	ws.Update("foo", true)
	ws.Update("foo", true)
	// TODO
}

func TestUpdateIncorrectThenCorrect(t *testing.T) {
	// Scheduled items should be empty.
	ws := reviewScheduler()
	ws.Update("foo", false)
	ws.Update("foo", true)

	items, _ := ws.ScheduleNow(-1)
	if len(items) > 0 {
		t.Log("expected items to be empty", items)
		t.Fail()
	}
}
