// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/lggruspe/polycloze/basedir"
)

type LanguageStats struct {
	// all-time
	Seen  int `json:"seen"`
	Total int `json:"total"`

	// today
	Learned  int `json:"learned"`
	Reviewed int `json:"reviewed"`

	// today
	Correct int `json:"correct"`
}

func queryInt(path, query string) (int, error) {
	var result int

	db, err := openDB(path)
	if err != nil {
		return result, err
	}
	defer db.Close()

	row := db.QueryRow(query)
	err = row.Scan(&result)
	return result, err
}

func countSeen(lang string) (int, error) {
	return queryInt(basedir.Review(lang), `select count(*) from review`)
}

// Total count of words in lang (given as three-letter code).
func countTotal(lang string) (int, error) {
	pattern := fmt.Sprintf("[a-z][a-z][a-z]-%s.db", lang)
	matches, _ := filepath.Glob(path.Join(basedir.DataDir, pattern))

	var max int
	for _, match := range matches {
		count, err := queryInt(match, `select count(*) from word`)
		if err != nil {
			return max, err
		}
		if count > max {
			max = count
		}
	}
	return max, nil
}

// New words learned today.
func countLearnedToday(lang string) (int, error) {
	query := `select count(*) from review where learned >= current_date`
	return queryInt(basedir.Review(lang), query)
}

// Number of words reviewed today, excluding new words.
func countReviewedToday(lang string) (int, error) {
	query := `
select count(*) from review where reviewed >= current_date
and learned < current_date
`
	return queryInt(basedir.Review(lang), query)
}

// Number of correct answers today.
func countCorrectToday(lang string) (int, error) {
	// NOTE assumes that 1 day is the smallest non-empty interval
	query := `select count(*) from review where reviewed >= current_date and correct`
	return queryInt(basedir.Review(lang), query)
}

func getLanguageStats(lang string) (*LanguageStats, error) {
	seen, err := countSeen(lang)
	if err != nil {
		return nil, err
	}

	total, err := countTotal(lang)
	if err != nil {
		return nil, err
	}

	learned, err := countLearnedToday(lang)
	if err != nil {
		return nil, err
	}

	reviewed, err := countReviewedToday(lang)
	if err != nil {
		return nil, err
	}

	correct, err := countCorrectToday(lang)
	if err != nil {
		return nil, err
	}

	return &LanguageStats{
		Seen:     seen,
		Total:    total,
		Learned:  learned,
		Reviewed: reviewed,
		Correct:  correct,
	}, nil
}
