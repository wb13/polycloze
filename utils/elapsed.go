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
