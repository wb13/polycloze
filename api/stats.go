// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"fmt"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
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

// If upgrade is non-empty, upgrades the database.
func queryInt(path, query string, upgrade ...bool) (int, error) {
	var result int

	db, err := database.Open(path)
	if err != nil {
		return 0, fmt.Errorf("could not open db (%v) for query (%v): %v", path, query, err)
	}
	defer db.Close()

	if len(upgrade) > 0 {
		if err := database.Upgrade(db); err != nil {
			return result, err
		}
	}

	row := db.QueryRow(query)
	err = row.Scan(&result)
	return result, err
}

func countSeen(l1, l2 string, userID int) (int, error) {
	return queryInt(basedir.Review(userID, l1, l2), `select count(*) from review`, true)
}

// Total count of words in course.
func countTotal(l1, l2 string) (int, error) {
	return queryInt(basedir.Course(l1, l2), `select count(*) from word`)
}

// New words learned today.
func countLearnedToday(l1, l2 string, userID int) (int, error) {
	query := `select count(*) from review where learned >= current_date`
	return queryInt(basedir.Review(userID, l1, l2), query, true)
}

// Number of words reviewed today, excluding new words.
func countReviewedToday(l1, l2 string, userID int) (int, error) {
	query := `
select count(*) from review where reviewed >= current_date
and learned < current_date
`
	return queryInt(basedir.Review(userID, l1, l2), query, true)
}

// Number of correct answers today.
func countCorrectToday(l1, l2 string, userID int) (int, error) {
	// NOTE assumes that 1 day is the smallest non-empty interval
	query := `select count(*) from review where reviewed >= current_date and correct`
	return queryInt(basedir.Review(userID, l1, l2), query, true)
}

func getCourseStats(l1, l2 string, userID int) (*CourseStats, error) {
	seen, err := countSeen(l1, l2, userID)
	if err != nil {
		return nil, err
	}

	total, err := countTotal(l1, l2)
	if err != nil {
		return nil, err
	}

	learned, err := countLearnedToday(l1, l2, userID)
	if err != nil {
		return nil, err
	}

	reviewed, err := countReviewedToday(l1, l2, userID)
	if err != nil {
		return nil, err
	}

	correct, err := countCorrectToday(l1, l2, userID)
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
