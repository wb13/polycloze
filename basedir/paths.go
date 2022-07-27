// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package basedir

import (
	"fmt"
	"path"
)

// Returns path to review database.
// l1 and l2: ISO 639-3 code
// TODO user param
func Review(l1, l2 string) string {
	return path.Join(StateDir, "reviews", "user", fmt.Sprintf("%s-%s.db", l1, l2))
}

// Returns path to log file.
// l1 and l2: ISO 639-3 code
// TODO user param
func Log(l1, l2 string) string {
	return path.Join(StateDir, "logs", "user", fmt.Sprintf("%s-%s.log", l1, l2))
}

// Returns path to database for course.
// l1 and l2 are ISO 639-3 codes.
func Course(l1, l2 string) string {
	return path.Join(DataDir, fmt.Sprintf("%s-%s.db", l1, l2))
}
