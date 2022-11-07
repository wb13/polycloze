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
