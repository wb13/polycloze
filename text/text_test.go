// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package text

import (
	"testing"
)

func TestCasefold(t *testing.T) {
	t.Parallel()

	a := Casefold("Foo")
	b := Casefold("fOo")
	c := Casefold("foO")

	if a != "foo" {
		t.Fatal("expected `a` to be 'foo':", a)
	}
	if b != "foo" {
		t.Fatal("expected `b` to be 'foo':", b)
	}
	if c != "foo" {
		t.Fatal("expected `c` to be 'foo':", c)
	}
}

func TestRemoveSoftHyphen(t *testing.T) {
	t.Parallel()

	examples := []string{
		Casefold("hell\u00ADo"),
		Casefold("hell­o"), // Contains soft-hyphen
	}

	for _, s := range examples {
		if s != "hello" {
			t.Fatal("expected `s` to be 'hello':", s)
		}
	}
}

func TestRemoveZeroWidthSpace(t *testing.T) {
	// Zero-width space at both ends should be removed, but not those in the middle.
	t.Parallel()

	examples := []string{
		Casefold("\u200Bhello"),
		Casefold("hello\u200B"),
		Casefold("\u200Bhello\u200B\u200B"),
	}

	for _, s := range examples {
		if s != "hello" {
			t.Fatal("expected `s` to be 'hello':", s)
		}
	}

	if Casefold("hell\u200Bo") == "hello" {
		t.Fatal("expected zero-width space in the middle to be left alone")
	}
}

func TestRemoveNoBreakSpace(t *testing.T) {
	// No-break space at both ends should be removed, but not those in the middle.
	t.Parallel()

	examples := []string{
		Casefold("\u00A0hello"),
		Casefold("hello\u00A0"),
		Casefold("\u00A0hello\u00A0\u00A0"),
	}

	for _, s := range examples {
		if s != "hello" {
			t.Fatal("expected `s` to be 'hello':", s)
		}
	}

	if Casefold("hell\u00A0o") == "hello" {
		t.Fatal("expected no-break space in the middle to be left alone")
	}
}
