// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/polycloze/polycloze/basedir"
	"github.com/polycloze/polycloze/database"
)

type Course struct {
	L1 Language `json:"l1"`
	L2 Language `json:"l2"`
}

// Checks if course exists.
func courseExists(l1, l2 string) bool {
	path := basedir.Course(l1, l2)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Gets user's active course.
// Similar to `getActiveCourse`, but takes the user ID instead of the user's
// database, and returns a Course value.
// Only use this if you don't plan to get other data from the user DB.
// Returns an error if the user hasn't set an active course.
func getUserActiveCourse(userID int) (Course, error) {
	// Open the user DB.
	db, err := database.OpenUserDB(basedir.UserData(userID))
	if err != nil {
		return Course{}, fmt.Errorf("failed to get active course: %w", err)
	}
	defer db.Close()

	// Get the course code.
	code, err := getActiveCourse(db)
	if err != nil {
		return Course{}, fmt.Errorf("failed to get active course: %w", err)
	}

	path := filepath.Join(basedir.DataDir, "courses", fmt.Sprintf("%v.db", code))
	course, err := getCourseInfo(path)
	if err != nil {
		return Course{}, fmt.Errorf("failed to get active course: %w", err)
	}
	return course, nil
}
