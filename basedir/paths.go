package basedir

import (
	"fmt"
	"path"
)

// Returns path to language database.
// language: ISO 630-3 code
func Language(language string) string {
	return path.Join(DataDir, "languages", fmt.Sprintf("%v.db", language))
}

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

// Returns path to translation database.
func Translation(l1, l2 string) string {
	if l1 == l2 {
		panic("invalid language pair")
	}
	if l2 < l1 {
		l1, l2 = l2, l1
	}
	pair := fmt.Sprintf("%s-%s.db", l1, l2)
	return path.Join(DataDir, "translations", pair)
}
