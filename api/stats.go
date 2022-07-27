// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"github.com/lggruspe/polycloze/basedir"
)

type CourseStats struct {
	// all-time
	Seen  int `json:"seen"`
	Total int `json:"total"`

	// today
	Learned  int `json:"learned"`
	Reviewed int `json:"reviewed"`
	Correct  int `json:"correct"`
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

func countSeen(l1, l2 string) (int, error) {
	return queryInt(basedir.Review(l1, l2), `select count(*) from review`)
}

// Total count of words in course.
func countTotal(l1, l2 string) (int, error) {
	return queryInt(basedir.Course(l1, l2), `select count(*) from word`)
}

// New words learned today.
func countLearnedToday(l1, l2 string) (int, error) {
	query := `select count(*) from review where learned >= current_date`
	return queryInt(basedir.Review(l1, l2), query)
}

// Number of words reviewed today, excluding new words.
func countReviewedToday(l1, l2 string) (int, error) {
	query := `
select count(*) from review where reviewed >= current_date
and learned < current_date
`
	return queryInt(basedir.Review(l1, l2), query)
}

// Number of correct answers today.
func countCorrectToday(l1, l2 string) (int, error) {
	// NOTE assumes that 1 day is the smallest non-empty interval
	query := `select count(*) from review where reviewed >= current_date and correct`
	return queryInt(basedir.Review(l1, l2), query)
}

func getCourseStats(l1, l2 string) (*CourseStats, error) {
	seen, err := countSeen(l1, l2)
	if err != nil {
		return nil, err
	}

	total, err := countTotal(l1, l2)
	if err != nil {
		return nil, err
	}

	learned, err := countLearnedToday(l1, l2)
	if err != nil {
		return nil, err
	}

	reviewed, err := countReviewedToday(l1, l2)
	if err != nil {
		return nil, err
	}

	correct, err := countCorrectToday(l1, l2)
	if err != nil {
		return nil, err
	}

	return &CourseStats{
		Seen:     seen,
		Total:    total,
		Learned:  learned,
		Reviewed: reviewed,
		Correct:  correct,
	}, nil
}
