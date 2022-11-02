// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package sessions

import (
	"testing"
)

// Generates session ID for testing.
// May generate duplicates.
func tid() string {
	id, err := generateID()
	if err != nil {
		panic(err)
	}
	return id
}

func TestCheckCSRFToken(t *testing.T) {
	t.Parallel()

	id := tid()
	token := CSRFToken(id)
	if !CheckCSRFToken(id, token) {
		t.Fatal("expected token to be valid for ID")
	}
}

func TestCSRFTokenIdNotEqual(t *testing.T) {
	t.Parallel()

	id := tid()
	token := CSRFToken(id)
	if id == token {
		t.Fatal("expected token and session ID to be different:", id, token)
	}
}
