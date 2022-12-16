// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package vocab_size

import (
	"testing"
	"time"

	"github.com/lggruspe/polycloze/review_scheduler"
	"github.com/lggruspe/polycloze/utils"
)

func TestVocabSizeEmptyRange(t *testing.T) {
	// Shouldn't crash, and result should be an empty slice.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	from := time.Now()
	series, err := VocabSize(db, from, from, time.Second)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(series) > 0 {
		t.Fatal("expected result to be empty:", series)
	}
}

func TestVocabSizeNoReviews(t *testing.T) {
	// Vocab size should be zero at any point in time.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	to := time.Now()
	from := to.AddDate(0, 0, -1)

	series, err := VocabSize(db, from, to, time.Hour)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(series) == 0 {
		t.Fatal("expected result to be non-empty:", series)
	}

	for _, metric := range series {
		if metric.Value != 0 {
			t.Fatal("expected vocab size to be 0:", metric.Value)
		}
	}
}

func TestVocabSizeInvalidResolution(t *testing.T) {
	// Panic if resolution is too high.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	defer func() { _ = recover() }()

	to := time.Now()
	from := to.AddDate(0, 0, -1)
	_, _ = VocabSize(db, from, to, time.Second-1) // Should panic.

	t.Fatal("did not panic")
}

func TestVocabSizeIncrease(t *testing.T) {
	// Successful review of previously unseen word should increase vocab size.
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

	// Get vocabulary size.
	series, err := VocabSize(db, from, to, 24*time.Hour)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(series) != 1 {
		t.Fatal("expected result to contain one partition:", series)
	}
}
