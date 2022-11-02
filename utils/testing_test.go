// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package utils

import (
	"testing"
)

func TestTestingDatabase(t *testing.T) {
	// TestingDatabase should run with no errors.
	t.Parallel()

	db := TestingDatabase()
	db.Close()
}
