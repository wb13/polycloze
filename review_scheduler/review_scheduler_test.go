// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package review_scheduler

import (
	"testing"
	"time"

	"github.com/lggruspe/polycloze/utils"
)

func TestSchedule(t *testing.T) {
	// Result should be empty with no errors.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	items, err := ScheduleReviewNow(db, 100)
	if err != nil {
		t.Fatal("expected err to be nil", err)
	}
	if len(items) > 0 {
		t.Error("expected items to be empty", items)
	}
}

func TestUpdate(t *testing.T) {
	// Only incorrect review needs to be reviewed.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	if err := UpdateReview(db, "foo", false); err != nil {
		t.Fatal("expected err to be nil", err)
	}
	if err := UpdateReview(db, "bar", true); err != nil {
		t.Fatal("expected err to be nil", err)
	}

	items, err := ScheduleReviewNow(db, 100)
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
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	items := []string{"foo", "bar", "baz"}
	for _, item := range items {
		if err := UpdateReview(db, item, true); err != nil {
			t.Fatal("expected err to be nil:", err)
		}
	}

	items, err := ScheduleReviewNow(db, -1)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if len(items) > 0 {
		t.Fatal("expected items to be empty", items)
	}
}

func TestUpdateIncorrectThenCorrect(t *testing.T) {
	// Scheduled items should be empty.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	if err := UpdateReview(db, "foo", false); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if err := UpdateReview(db, "foo", true); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	items, _ := ScheduleReviewNow(db, -1)
	if len(items) > 0 {
		t.Log("expected items to be empty", items)
		t.Fail()
	}
}

func TestUpdateSuccessfulReviewDoesNotDecreaseIntervalSize(t *testing.T) {
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	query := func() time.Duration {
		row := db.QueryRow(`select interval from review`)
		var interval time.Duration
		if err := row.Scan(&interval); err != nil {
			t.Fatal("expected err to be nil:", err)
		}
		return interval * time.Second
	}

	now := time.Now().UTC()

	if err := UpdateReviewAt(db, "foo", true, now); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	before := query()

	now = now.Add(3 * 24 * time.Hour)
	if err := UpdateReviewAt(db, "foo", true, now); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	after := query()

	if before > after {
		t.Fatal("expected sequence of successful reviews to have non-decreasing intervals", before, after)
	}
}

func TestCase(t *testing.T) {
	// Items shouldn't be case-folded.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	if err := UpdateReview(db, "Foo", false); err != nil {
		t.Fatal("expected nil err", err)
	}

	items, err := ScheduleReviewNow(db, 100)
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

func TestReviewTimestampType(t *testing.T) {
	// Timestamps should be stored as integers (UNIX timestamps).
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	if err := UpdateReview(db, "foo", true); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	query := `
		SELECT learned, typeof(learned), reviewed, typeof(reviewed), due, typeof(due)
		FROM review
	`

	rows, err := db.Query(query)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var learned, reviewed, due int
		var tLearned, tReviewed, tDue string

		if err := rows.Scan(&learned, &tLearned, &reviewed, &tReviewed, &due, &tDue); err != nil {
			t.Fatal("expected err to be nil:", err)
		}

		if tLearned != "integer" {
			t.Fatal("expected typeof(learned) to be 'integer':", tLearned)
		}
		if tReviewed != "integer" {
			t.Fatal("expected typeof(reviewed) to be 'integer':", tReviewed)
		}
		if tDue != "integer" {
			t.Fatal("expected typeof(due) to be 'integer':", tDue)
		}

		// Check if learned, reviewed and due have more than 4 digits.
		// `cast(current_timestamp as integer)` returns a 4 digit number.
		if learned < 10000 {
			t.Fatal("expected learned to be a UNIX timestamp:", learned)
		}
		if reviewed < 10000 {
			t.Fatal("expected reviewed to be a UNIX timestamp:", reviewed)
		}
		if due < 10000 {
			t.Fatal("expected due to be a UNIX timestamp:", due)
		}
	}
}
