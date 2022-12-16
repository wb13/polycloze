// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package history

import (
	"database/sql"
	"fmt"
	"time"
)

// Represents a record in the review history.
type Review struct {
	Word           string
	Reviewed       time.Time
	IntervalBefore time.Duration // Negative if null
	IntervalAfter  time.Duration
}

// Summary of reviews within an interval of time.
type Summary struct {
	From time.Time
	To   time.Time

	Unimproved   int
	Learned      int
	Forgotten    int
	Crammed      int
	Strengthened int
}

// Summarizes review activity during the given range.
// The range gets partitioned into intervals of length `step`.
// The result contains a summary for each interval.
// Only supports up to second precision for `step`.
func Summarize(db *sql.DB, from, to time.Time, step time.Duration) ([]Summary, error) {
	if step < time.Second {
		panic("only supports up to second precision")
	}

	// Initialize return value.
	var summaries []Summary
	for current := from; current.Before(to); current = current.Add(step) {
		summaries = append(summaries, Summary{
			From: current,
			To:   current.Add(step),
		})
	}

	// Compute summaries.
	query := `
		SELECT (reviewed - @from)/@step, coalesce(interval_before, 0), interval_after
		FROM history
		WHERE reviewed >= @from AND reviewed < @to
		ORDER BY reviewed ASC
	`
	rows, err := db.Query(
		query,
		sql.Named("from", from.Unix()),
		sql.Named("to", to.Unix()),
		sql.Named("step", step/time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to summarize review history: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i, intervalBefore, intervalAfter int64
		if err := rows.Scan(&i, &intervalBefore, &intervalAfter); err != nil {
			return nil, fmt.Errorf("failed to summarize review history: %v", err)
		}

		if intervalBefore <= 0 && intervalAfter <= 0 {
			summaries[i].Unimproved++
		} else if intervalBefore <= 0 && intervalAfter > 0 {
			summaries[i].Learned++
		} else if intervalBefore > intervalAfter {
			summaries[i].Forgotten++
		} else if intervalBefore == intervalAfter {
			summaries[i].Crammed++
		} else if intervalBefore < intervalAfter {
			summaries[i].Strengthened++
		}
	}
	return summaries, nil
}

// Returns list of reviews in given range.
// The list starts with newer reviews.
// Specify a negative `limit` to return an unbounded number of `Review`s.
func Get(db *sql.DB, from, to time.Time, step time.Duration, limit int) ([]Review, error) {
	query := `
		SELECT word, reviewed, coalesce(interval_before, -1), interval_after
		FROM history
		WHERE reviewed >= ? AND reviewed < ?
		ORDER BY reviewed DESC
		LIMIT ?
	`

	rows, err := db.Query(query, from.Unix(), to.Unix(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get past reviews: %v", err)
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		if limit >= 0 && len(reviews) >= limit {
			break
		}

		var reviewed int64
		var intervalBefore, intervalAfter time.Duration
		var review Review

		err = rows.Scan(&review.Word, &reviewed, &intervalBefore, &intervalAfter)
		if err != nil {
			return nil, fmt.Errorf("failed to get past reviews: %v", err)
		}

		review.Reviewed = time.Unix(reviewed, 0)
		review.IntervalBefore = intervalBefore * time.Hour
		review.IntervalAfter = intervalAfter * time.Hour
		reviews = append(reviews, review)
	}
	return reviews, nil
}
