// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package review_scheduler

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lggruspe/polycloze/utils"
)

func queryIntervals(tx *sql.Tx) []time.Duration {
	query := `SELECT interval FROM interval ORDER BY interval ASC`
	rows, err := tx.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var intervals []time.Duration
	for rows.Next() {
		var interval time.Duration
		if err := rows.Scan(&interval); err != nil {
			return nil
		}
		intervals = append(intervals, interval*time.Hour)
	}
	return intervals
}

func TestInsertInterval(t *testing.T) {
	// Intervals should be stored as number of hours.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if err := insertInterval(tx, time.Hour); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	var interval int64
	query := `SELECT max(interval) FROM interval`
	if err := db.QueryRow(query).Scan(&interval); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if interval != 1 {
		t.Fatal("expected `interval` to be equal to 1:", interval)
	}
}

func TestMaxInterval(t *testing.T) {
	// Should return 0 if there are no intervals in the database.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	max, err := maxInterval(tx)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if max != 0 {
		t.Fatal("expected max interval to be 0:", err)
	}
}

func TestShortenIntervalLessThanDay(t *testing.T) {
	// shortenInterval shouldn't do anything if the interval is shorter than a day.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Insert some intervals
	for i := 0; i < 25; i++ {
		if err := insertInterval(tx, time.Duration(i)*time.Hour); err != nil {
			t.Fatal("expected err to be nil:", err)
		}
	}

	// Try to shorten all intervals.
	for i := 0; i < 25; i++ {
		if err := shortenInterval(tx, time.Duration(i)*time.Hour); err != nil {
			t.Fatal("expected err to be nil:", err)
		}
	}

	// Check to make sure none of the intervals were modified
	intervals := queryIntervals(tx)
	if n := len(intervals); n != 25 {
		t.Fatal("expected number of intervals to not change:", intervals, n)
	}
	for i, interval := range intervals {
		if hours := int64(interval.Hours()); int64(i) != hours {
			t.Fatal("expected intervals to not change:", i, hours)
		}
	}
}

func TestShortenIntervalNotExisting(t *testing.T) {
	// If the replacement interval doesn't exist, simply change the value of the old interval.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	intervals := []time.Duration{
		0,
		time.Duration(24) * time.Hour,
		time.Duration(48) * time.Hour,
		time.Duration(96) * time.Hour,
	}
	for _, interval := range intervals {
		if err := insertInterval(tx, interval); err != nil {
			t.Fatal("expected err to be nil:", err)
		}
	}

	if err := shortenInterval(tx, intervals[2]); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	result := queryIntervals(tx)

	if len(intervals) != len(result) || len(result) != 4 {
		t.Fatal("size of result different than expected (4):", result)
	}

	// All intervals should be the same, except for 48.
	// It should have been shortened to 36 by now.
	for i, interval := range intervals {
		if i != 2 && interval != result[i] {
			t.Fatal("expected non-shortened intervals to remain the same:", interval, result[i])
		}
	}

	if intervals[2] == result[2] {
		t.Fatal("expected shortened interval to change:", intervals[2], result[2])
	}

	if expected := time.Duration(36) * time.Hour; result[2] != expected {
		t.Fatal("expected result to be 36 hours:", result[2], expected)
	}
}

func TestShortenIntervalExisting(t *testing.T) {
	// If the replacement interval exists already, the old interval should be deleted.
	// This happens when `replacement = old interval - 1` already exists.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Insert some intervals.
	intervals := []time.Duration{
		0,
		time.Duration(24) * time.Hour,
		time.Duration(25) * time.Hour,
		time.Duration(26) * time.Hour,
		time.Duration(27) * time.Hour,
	}
	for _, interval := range intervals {
		if err := insertInterval(tx, interval); err != nil {
			t.Fatal("expected err to be nil:", err)
		}
	}

	// Shorten 26 hour interval.
	if err := shortenInterval(tx, 26*time.Hour); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// The 26 hour interval should be deleted
	intervals = queryIntervals(tx)
	if n := len(intervals); n != 4 {
		t.Fatal("expected number of intervals to decrease by one:", intervals, n)
	}

	expected := []time.Duration{
		0,
		time.Duration(24) * time.Hour,
		time.Duration(25) * time.Hour,
		time.Duration(27) * time.Hour,
	}

	if len(intervals) != len(expected) {
		t.Fatal("result different than expected:", intervals, expected)
	}

	for i, interval := range intervals {
		if interval != expected[i] {
			t.Fatal("result different than expected:", intervals, expected)
		}
	}
}

func TestLengthenIntervalLessThanDay(t *testing.T) {
	// Ther should be no changes if the interval is shorter than 1 day.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Insert some intervals.
	for i := 0; i < 25; i++ {
		if err := insertInterval(tx, time.Duration(i)*time.Hour); err != nil {
			t.Fatal("expected err to be nil:", err)
		}
	}

	// Try to lengthen all intervals.
	for i := 0; i < 25; i++ {
		if err := lengthenInterval(tx, time.Duration(i)*time.Hour); err != nil {
			t.Fatal("expected err to be nil:", err)
		}
	}

	// Check to make sure none of the intervals were modified
	intervals := queryIntervals(tx)
	if n := len(intervals); n != 25 {
		t.Fatal("expected number of intervals to not change:", intervals, n)
	}
	for i, interval := range intervals {
		if hours := int64(interval.Hours()); int64(i) != hours {
			t.Fatal("expected intervals to not change:", i, hours)
		}
	}
}

func TestLengthenIntervalExisting(t *testing.T) {
	// If the replacement interval already exists, the old interval should be deleted.
	// This happens when `replacement = old interval + 1` already exists.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Insert some intervals
	intervals := []time.Duration{
		0,
		time.Duration(24) * time.Hour,
		time.Duration(25) * time.Hour,
		time.Duration(26) * time.Hour,
		time.Duration(27) * time.Hour,
	}
	for _, interval := range intervals {
		if err := insertInterval(tx, interval); err != nil {
			t.Fatal("expected err to be nil:", err)
		}
	}

	// Lengthen 25 hour interval.
	if err := lengthenInterval(tx, 25*time.Hour); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// The 25 hour interval should be deleted
	intervals = queryIntervals(tx)
	if n := len(intervals); n != 4 {
		t.Fatal("expected number of intervals to decrease by one:", intervals, n)
	}

	expected := []time.Duration{
		0,
		time.Duration(24) * time.Hour,
		time.Duration(26) * time.Hour,
		time.Duration(27) * time.Hour,
	}

	if len(intervals) != len(expected) {
		t.Fatal("result different than expected:", intervals, expected)
	}

	for i, interval := range intervals {
		if interval != expected[i] {
			t.Fatal("result different than expected:", intervals, expected)
		}
	}
}

func TestLengthenIntervalNotExisting(t *testing.T) {
	// If the replacement interval doesn't exist, simply change the value of the old interval.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Insert some intervals
	intervals := []time.Duration{
		0,
		time.Duration(24) * time.Hour,
		time.Duration(48) * time.Hour,
		time.Duration(96) * time.Hour,
	}
	for _, interval := range intervals {
		if err := insertInterval(tx, interval); err != nil {
			t.Fatal("expected err to be nil:", err)
		}
	}

	// Lengthen the 48h interval.
	if err := lengthenInterval(tx, intervals[2]); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Check resulting intervals.
	result := queryIntervals(tx)

	if len(intervals) != len(result) || len(result) != 4 {
		t.Fatal("size of result different than expected (4):", result)
	}

	// All intervals should be the same, except for 48, which should now be 72.
	for i, interval := range intervals {
		if i != 2 && interval != result[i] {
			t.Fatal("expected non-lengthened intervals to remain the same:", interval, result[i])
		}
	}

	if intervals[2] == result[2] {
		t.Fatal("expected lengthened interval to change:", intervals[2], result[2])
	}

	if expected := time.Duration(72) * time.Hour; result[2] != expected {
		t.Fatal("expected result to be 72 hours:", result[2], expected)
	}
}
