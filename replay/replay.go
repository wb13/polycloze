// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package replay

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/word_scheduler"
)

// Checks if `io.Reader` has reached the EOF.
// Assumes the reader will no longer be used.
func isEOF(r io.Reader) bool {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		return false
	}
	return scanner.Err() == nil
}

// Imports review data from CSV file.
func Replay[T database.Querier](q T, r io.Reader) error {
	// TODO don't allow this operation if there are existing reviews in the DB.
	reader := NewReviewReader(csv.NewReader(r))

	// Ignore first error (it may be a header row), but don't ignore further
	// errors.
	if review, err := reader.ReadReview(); err == nil {
		if err := word_scheduler.UpdateWordAt(
			q,
			review.Word,
			review.Correct,
			review.Reviewed,
		); err != nil {
			return fmt.Errorf("failed to import review: %v", err)
		}
	}

	var review ReviewEvent
	var err error
	for {
		review, err = reader.ReadReview()
		if err != nil {
			break
		}
		if err := word_scheduler.UpdateWordAt(
			q,
			review.Word,
			review.Correct,
			review.Reviewed,
		); err != nil {
			return fmt.Errorf("failed to import review: %v", err)
		}
	}

	if !isEOF(r) {
		return fmt.Errorf("failed to import review: %v", err)
	}
	return nil
}

func ReplayFile[T database.Querier](q T, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to import reviews from file: %v", err)
	}
	if err := Replay(q, f); err != nil {
		return fmt.Errorf("failed to import reviews from file: %v", err)
	}
	return nil
}
