// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"os"

	"github.com/lggruspe/polycloze/basedir"
)

// Checks if course exists.
func courseExists(l1, l2 string) bool {
	path := basedir.Course(l1, l2)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
