// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package database

import (
	"testing"

	"github.com/lggruspe/polycloze/basedir"
)

func TestSessionOpenClose(t *testing.T) {
	t.Parallel()

	db, err := New(":memory:")
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	s, err := NewSession(db, basedir.Course("eng", "spa"))
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if err := s.Close(); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
}
