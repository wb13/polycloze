package review_scheduler

import (
	"testing"
	"time"

	"github.com/lggruspe/polycloze/database"
)

// Returns ReviewScheduler for testing.
func reviewScheduler() *database.Session {
	db, _ := New(":memory:")
	s, _ := database.NewSession(db, "", "", "")
	return s
}

func TestSchedule(t *testing.T) {
	// Result should be empty with no errors.
	s := reviewScheduler()

	items, err := ScheduleReviewNow(s, 100)

	if err != nil {
		t.Fatal("expected err to be nil", err)
	}
	if len(items) > 0 {
		t.Error("expected items to be empty", items)
	}
}

func TestUpdate(t *testing.T) {
	// Only incorrect review needs to be reviewed.
	s := reviewScheduler()

	if err := UpdateReview(s, "foo", false); err != nil {
		t.Fatal("expected err to be nil", err)
	}
	if err := UpdateReview(s, "bar", true); err != nil {
		t.Fatal("expected err to be nil", err)
	}

	items, err := ScheduleReviewNow(s, 100)
	if err != nil {
		t.Fatal("expected err to be nil", err)
	}

	if len(items) != 1 {
		t.Log("expected different number of results", items)
	}
	if items[0] != "foo" {
		t.Error("expected scheduled items to contain \"foo\"", items[0])
	}
}

func TestUpdateRecentlyAnsweredItemDoesntGetScheduled(t *testing.T) {
	s := reviewScheduler()
	items := []string{"foo", "bar", "baz"}
	for _, item := range items {
		UpdateReview(s, item, true)
	}

	items, err := ScheduleReviewNow(s, -1)
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
	s := reviewScheduler()
	UpdateReview(s, "foo", false)
	UpdateReview(s, "foo", true)

	items, _ := ScheduleReviewNow(s, -1)
	if len(items) > 0 {
		t.Log("expected items to be empty", items)
		t.Fail()
	}
}

func TestUpdateSuccessfulReviewDoesNotDecreaseIntervalSize(t *testing.T) {
	s := reviewScheduler()
	UpdateReview(s, "foo", true)
	UpdateReview(s, "foo", true)
	UpdateReview(s, "foo", true)

	query := `SELECT interval FROM review ORDER BY reviewed ASC`
	rows, err := s.Query(query)
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

func TestCase(t *testing.T) {
	// Items shouldn't be case-folded.
	s := reviewScheduler()

	if err := UpdateReview(s, "Foo", false); err != nil {
		t.Fatal("expected nil err", err)
	}

	items, err := ScheduleReviewNow(s, 100)
	if err != nil {
		t.Fatal("expected nil err", err)
	}

	if len(items) != 1 {
		t.Fatal("expected different number of results", items)
	}

	if items[0] != "Foo" {
		t.Error("expected \"Foo\"")
	}
}
