// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package text

import (
	"strings"

	"golang.org/x/text/cases"
)

var caser = cases.Fold()

// Casefolds string and removes soft-hyphens.
func Casefold(s string) string {
	s = strings.ReplaceAll(s, "\u00AD", "")
	return caser.String(s)
}
