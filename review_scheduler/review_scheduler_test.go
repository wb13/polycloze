//go:build sqlite_math_functions

package review_scheduler

import (
	"database/sql"
	"testing"
	"time"
)

// Returns ReviewScheduler for testing.
func reviewScheduler() *sql.DB {
	db, _ := New(":memory:")
	return db
}

func TestSchedule(t *testing.T) {
	// Result should be empty with no errors.
	db := reviewScheduler()

	items, err := ScheduleReviewNow(db, 100)

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
	db := reviewScheduler()

	if err := UpdateReview(db, "foo", false); err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}
	if err := UpdateReview(db, "bar", true); err != nil {
		t.Log("expected err to be nil", err)
		t.Fail()
	}

	items, err := ScheduleReviewNow(db, 100)
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
	db := reviewScheduler()
	items := []string{"foo", "bar", "baz"}
	for _, item := range items {
		UpdateReview(db, item, true)
	}

	items, err := ScheduleReviewNow(db, -1)
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
	db := reviewScheduler()
	UpdateReview(db, "foo", false)
	UpdateReview(db, "foo", true)

	items, _ := ScheduleReviewNow(db, -1)
	if len(items) > 0 {
		t.Log("expected items to be empty", items)
		t.Fail()
	}
}

func TestUpdateSuccessfulReviewDoesNotDecreaseIntervalSize(t *testing.T) {
	db := reviewScheduler()
	UpdateReview(db, "foo", true)
	UpdateReview(db, "foo", true)
	UpdateReview(db, "foo", true)

	query := `SELECT interval FROM review ORDER BY id ASC`
	rows, err := db.Query(query)
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
