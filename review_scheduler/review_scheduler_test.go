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

	query := func() time.Duration {
		row := s.QueryRow(`select interval from review`)
		var interval time.Duration
		if err := row.Scan(&interval); err != nil {
			t.Fatal("expected err to be nil:", err)
		}
		return interval * time.Second
	}

	now := time.Now().UTC()

	if err := UpdateReviewAt(s, "foo", true, now); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	before := query()

	now = now.Add(3 * 24 * time.Hour)
	if err := UpdateReviewAt(s, "foo", true, now); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	after := query()

	if before > after {
		t.Fatal("expected sequence of successful reviews to have non-decreasing intervals", before, after)
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
