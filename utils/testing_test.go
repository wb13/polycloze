// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package utils

import (
	"testing"
)

func TestTestingDatabase(_ *testing.T) {
	// TestingDatabase should run with no errors.
	db := TestingDatabase()
	db.Close()
}
