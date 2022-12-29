// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package replay

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"
)

func testReader(s string) *ReviewReader {
	return NewReviewReader(csv.NewReader(strings.NewReader(s)))
}

func writeCSV(reviews []ReviewEvent) string {
	b := new(strings.Builder)
	w := NewReviewWriter(csv.NewWriter(b))

	for _, r := range reviews {
		if err := w.WriteReview(r); err != nil {
			panic(err)
		}
	}
	return b.String()
}

func TestReadReviewWord(t *testing.T) {
	t.Parallel()

	r := testReader(`word,reviewed,correct
foo,0,1
bar,0,1
`)
	// Skip header.
	_, _ = r.ReadReview()

	e, err := r.ReadReview()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if e.Word != "foo" {
		t.Fatal("expected record.Word to be 'foo':", e.Word)
	}

	e, err = r.ReadReview()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if e.Word != "bar" {
		t.Fatal("expected record.Word to be 'bar':", e.Word)
	}

	_, err = r.ReadReview()
	if err == nil {
		t.Fatal("expected no records after last row")
	}
}

func TestReadReviewCorrect(t *testing.T) {
	t.Parallel()

	r := testReader(`word,reviewed,correct
foo,0,0
bar,0,1
baz,0,2
`)

	// Skip header.
	_, _ = r.ReadReview()

	// Parse 0.
	e, err := r.ReadReview()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if e.Word != "foo" {
		t.Fatal("expected record.Word to be 'foo':", e.Word)
	}
	if e.Correct {
		t.Fatal("expected record.Correct to be false:", e.Correct)
	}

	// Parse 1.
	e, err = r.ReadReview()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if e.Word != "bar" {
		t.Fatal("expected record.Word to be 'bar':", e.Word)
	}
	if !e.Correct {
		t.Fatal("expected record.Correct to be false:", e.Correct)
	}

	// Parse invalid correct.
	if _, err := r.ReadReview(); err == nil {
		t.Fatal("expected invalid correct value to cause error")
	}
}

func TestReadReviewReviewed(t *testing.T) {
	t.Parallel()

	now := time.Now()
	reviews := []ReviewEvent{
		{
			Word:     "foo",
			Reviewed: time.Unix(0, 0),
			Correct:  false,
		},
		{
			Word:     "bar",
			Reviewed: now,
			Correct:  true,
		},
	}

	r := testReader(writeCSV(reviews))

	// Read first review.
	e, err := r.ReadReview()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if e.Word != "foo" {
		t.Fatal("expected record.Word to be 'foo':", e.Word)
	}
	if e.Reviewed != reviews[0].Reviewed {
		t.Fatal(
			"expected record.Reviewed to be the same:",
			e.Reviewed,
			reviews[0].Reviewed,
		)
	}

	// Read second review.
	e, err = r.ReadReview()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if e.Word != "bar" {
		t.Fatal("expected record.Word to be 'bar':", e.Word)
	}
	if a, b := e.Reviewed.Unix(), reviews[1].Reviewed.Unix(); a != b {
		t.Fatal("expected record.Reviewed to be the same:", a, b)
	}
}
