// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"fmt"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/word_scheduler"
)

// Saves review results to the db.
// Returns an error if it fails to save one or more of the review results.
// The caller may choose to ignore the error.
func saveReviewResults[T database.Querier](q T, reviews []ReviewResult) error {
	var err error
	for _, review := range reviews {
		_err := word_scheduler.UpdateWord(q, review.Word, review.Correct)
		if _err != nil {
			err = _err
		}
	}

	if err != nil {
		return fmt.Errorf("failed to save some reviews: %v", err)
	}
	return err
}
