// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package utils

import (
	"fmt"
	"time"
)

func Elapsed(description string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("[%s] took %v\n", description, time.Since(start))
	}
}
