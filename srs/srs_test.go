//go:build sqlite_math_functions

package srs

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestInitScheduler(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	rs, err := InitReviewScheduler(db)

	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}

	if rs.db == nil {
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
	rs, _ := InitReviewScheduler(db)
	return rs
}

func TestSchedule(t *testing.T) {
	// Result should be empty with no errors.
	rs := reviewScheduler()

	items, err := rs.ScheduleNow(100)

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
	rs := reviewScheduler()

	if err := rs.Update("foo", false); err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
	if err := rs.Update("bar", true); err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}

	items, err := rs.ScheduleNow(100)
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
	rs := reviewScheduler()
	items := []string{"foo", "bar", "baz"}
	for _, item := range items {
		rs.Update(item, true)
	}

	items, err := rs.ScheduleNow(-1)
	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
	if len(items) > 0 {
		t.Log("expected items to be empty", items)
		t.Fail()
	}
}

func TestUpdateIncorrectThenCorrect(t *testing.T) {
	// Scheduled items should be empty.
	rs := reviewScheduler()
	rs.Update("foo", false)
	rs.Update("foo", true)

	items, _ := rs.ScheduleNow(-1)
	if len(items) > 0 {
		t.Log("expected items to be empty", items)
		t.Fail()
	}
}

func TestUpdateSuccessfulReviewDoesNotDecreaseIntervalSize(t *testing.T) {
	rs := reviewScheduler()
	rs.Update("foo", true)
	rs.Update("foo", true)
	rs.Update("foo", true)

	query := `SELECT interval FROM Review ORDER BY id ASC`
	rows, err := rs.db.Query(query)
	if err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
	defer rows.Close()

	var intervals []time.Duration
	for rows.Next() {
		var interval time.Duration
		if err := rows.Scan(&interval); err != nil {
			t.Log("expected err to be nil", err)
			t.Fail()
		}
		intervals = append(intervals, interval)
	}

	for i := 1; i < len(intervals); i++ {
		if intervals[i-1] > intervals[i] {
			t.Log("expected sequence of successful reviews to have non-decreasing intervals", intervals)
			t.Fail()
		}
	}
}
