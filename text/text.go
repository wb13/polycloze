// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package text

import (
	"golang.org/x/text/cases"
)

var caser = cases.Fold()

func Casefold(s string) string {
	return caser.String(s)
}
