// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package text

import (
	"strings"

	"golang.org/x/text/cases"
)

const (
	softHyphen     = "\u00AD"
	zeroWidthSpace = "\u200B"
	noBreakSpace   = "\u00A0"
)

var caser = cases.Fold()

// Casefolds string and removes some unnedeed characters.
func Casefold(s string) string {
	// NOTE This operation is also performed in `python/scripts/word.py`, so any
	// changes here should be reflected there as well.
	s = strings.ReplaceAll(s, softHyphen, "")

	for strings.HasPrefix(s, zeroWidthSpace) {
		s = strings.TrimPrefix(s, zeroWidthSpace)
	}
	for strings.HasSuffix(s, zeroWidthSpace) {
		s = strings.TrimSuffix(s, zeroWidthSpace)
	}

	for strings.HasPrefix(s, noBreakSpace) {
		s = strings.TrimPrefix(s, noBreakSpace)
	}
	for strings.HasSuffix(s, noBreakSpace) {
		s = strings.TrimSuffix(s, noBreakSpace)
	}

	return caser.String(s)
}
