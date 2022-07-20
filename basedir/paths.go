package basedir

import (
	"fmt"
	"path"
)

// Returns path to review database.
// language: ISO 639-3 code
// TODO user param
func Review(language string) string {
	return path.Join(StateDir, "reviews", "user", fmt.Sprintf("%v.db", language))
}

// Returns path to log file.
// language: ISO 639-3 code
// TODO user param
func Log(language string) string {
	return path.Join(StateDir, "logs", "user", fmt.Sprintf("%v.log", language))
}

// Returns path to database for course.
// l1 and l2 are ISO 639-3 codes.
func Course(l1, l2 string) string {
	return path.Join(DataDir, fmt.Sprintf("%s-%s.db", l1, l2))
}
