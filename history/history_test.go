// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package history

import (
	"testing"
	"time"

	"github.com/polycloze/polycloze/review_scheduler"
	"github.com/polycloze/polycloze/utils"
)

func TestSummarizeEmptyRange(t *testing.T) {
	// Result should be an empty slice.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	from := time.Now()
	result, err := Summarize(db, from, from, time.Second)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(result) > 0 {
		t.Fatal("expected result to be empty:", result)
	}
}

func TestSummarizeNoReviews(t *testing.T) {
	// Summary should be 0 for all intervals.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	to := time.Now()
	from := to.AddDate(0, 0, -1)

	result, err := Summarize(db, from, to, time.Hour)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(result) == 0 {
		t.Fatal("expected result to be non-empty:", result)
	}

	for _, summary := range result {
		ok := summary.Unimproved == 0 &&
			summary.Learned == 0 &&
			summary.Forgotten == 0 &&
			summary.Crammed == 0 &&
			summary.Strengthened == 0
		if !ok {
			t.Fatal("expected summary values to be 0:", summary)
		}
	}
}

func TestSummarizeInvalidResolution(t *testing.T) {
	// Panic if resolution is too high.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	defer func() { _ = recover() }()

	to := time.Now()
	from := to.AddDate(0, 0, -1)

	// Panics.
	_, _ = Summarize(db, from, to, time.Second-1)

	t.Fatal("did not panic")
}

func TestSummarizeLearn(t *testing.T) {
	// Some of the summary values should be non-zero.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	to := time.Now()
	from := to.AddDate(0, 0, -1)

	// Review a word.
	at := from.Add(time.Hour)
	if err := review_scheduler.UpdateReviewAt(db, "foo", true, at); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Get summary.
	result, err := Summarize(db, from, to, 24*time.Hour)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(result) != 1 {
		t.Fatal("expected result to contain one partition:", result)
	}

	// Check summary values.
	summary := result[0]
	ok := summary.Unimproved == 0 &&
		summary.Learned == 1 && // for reviewing previously unseen word
		summary.Forgotten == 0 &&
		summary.Crammed == 0 &&
		summary.Strengthened == 0
	if !ok {
		t.Fatal("incorrect summary:", summary)
	}
}

func TestSummarizeCram(t *testing.T) {
	// Some of the summary values should be non-zero.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	to := time.Now()
	from := to.AddDate(0, 0, -1)

	// Review a word.
	at := from.Add(time.Hour)
	if err := review_scheduler.UpdateReviewAt(db, "foo", true, at); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Cram.
	at = at.Add(time.Hour)
	if err := review_scheduler.UpdateReviewAt(db, "foo", true, at); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Get summary.
	result, err := Summarize(db, from, to, 24*time.Hour)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(result) != 1 {
		t.Fatal("expected result to contain one partition:", result)
	}

	// Check summary values.
	summary := result[0]
	ok := summary.Unimproved == 0 &&
		summary.Learned == 1 && // 1 for initial
		summary.Forgotten == 0 &&
		summary.Crammed == 1 && // 1 for crammed review
		summary.Strengthened == 0
	if !ok {
		t.Fatal("incorrect summary:", summary)
	}
}

func TestSummarizeUnimproved(t *testing.T) {
	// Some of the summary values should be non-zero.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	to := time.Now()
	from := to.AddDate(0, 0, -1)

	// Review a word.
	at := from.Add(time.Hour)
	if err := review_scheduler.UpdateReviewAt(db, "foo", false, at); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Get summary.
	result, err := Summarize(db, from, to, 24*time.Hour)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(result) != 1 {
		t.Fatal("expected result to contain one partition:", result)
	}

	// Check summary values.
	summary := result[0]
	ok := summary.Unimproved == 1 && // for failed review
		summary.Learned == 0 &&
		summary.Forgotten == 0 &&
		summary.Crammed == 0 &&
		summary.Strengthened == 0
	if !ok {
		t.Fatal("incorrect summary:", summary)
	}
}

func TestSummarizeForgotten(t *testing.T) {
	// Some of the summary values should be non-zero.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	to := time.Now()
	from := to.AddDate(0, 0, -1)

	// Review a word.
	at := from.Add(time.Hour)
	if err := review_scheduler.UpdateReviewAt(db, "foo", true, at); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Forget the word.
	at = from.Add(time.Hour)
	if err := review_scheduler.UpdateReviewAt(db, "foo", false, at); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Get summary.
	result, err := Summarize(db, from, to, 24*time.Hour)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(result) != 1 {
		t.Fatal("expected result to contain one partition:", result)
	}

	// Check summary values.
	summary := result[0]
	ok := summary.Unimproved == 0 &&
		summary.Learned == 1 && // for initial review
		summary.Forgotten == 1 && // for second review
		summary.Crammed == 0 &&
		summary.Strengthened == 0
	if !ok {
		t.Fatal("incorrect summary:", summary)
	}
}

func TestSummarizeStrengthened(t *testing.T) {
	// Some of the summary values should be non-zero.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	to := time.Now()
	from := to.AddDate(0, 0, -2)

	// Review a word.
	at := from.Add(time.Hour)
	if err := review_scheduler.UpdateReviewAt(db, "foo", true, at); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Strengthen the word after the due date.
	at = from.Add(25 * time.Hour)
	if err := review_scheduler.UpdateReviewAt(db, "foo", true, at); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Get summary.
	result, err := Summarize(db, from, to, 48*time.Hour)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(result) != 1 {
		t.Fatal("expected result to contain one partition:", result)
	}

	// Check summary values.
	summary := result[0]
	ok := summary.Unimproved == 0 &&
		summary.Learned == 1 && // for initial review
		summary.Forgotten == 0 &&
		summary.Crammed == 0 &&
		summary.Strengthened == 1 // for second review
	if !ok {
		t.Fatal("incorrect summary:", summary)
	}
}
