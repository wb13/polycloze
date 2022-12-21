// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"fmt"
	"os"
	"strings"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
)

// Checks if course exists.
func courseExists(l1, l2 string) bool {
	path := basedir.Course(l1, l2)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Gets user's active course.
// Similar to `getActiveCourse`, but takes the user ID instead of the user's
// database.
// Also returns the L1 and L2 codes separately rather than a single course
// code.
// Only use this if you don't plan to get other data from the user DB.
// Returns an error if the user hasn't set an active course.
func getUserActiveCourse(userID int) (string, string, error) {
	// Open the user DB.
	db, err := database.OpenUserDB(basedir.UserData(userID))
	if err != nil {
		return "", "", fmt.Errorf("failed to get active course: %v", err)
	}
	defer db.Close()

	// Get the course code.
	code, err := getActiveCourse(db)
	if err != nil {
		return "", "", fmt.Errorf("failed to get active course: %v", err)
	}

	// Get L1 and L2 codes.
	l1Code, l2Code, found := strings.Cut(code, "-")
	if !found {
		// This shouldn't happen unless a course gets deleted or a course code
		// contains several dashes (not supported yet).
		return "", "", fmt.Errorf(
			"failed to get active course: unsupported course %v",
			code,
		)
	}
	return l1Code, l2Code, nil
}
