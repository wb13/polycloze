// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// For computing Wilson score intervals
package wilson

import (
	"math"
)

// Computes a boundary point of a Wilson score interval.
// See https://www.itl.nist.gov/div898/handbook/prc/section2/prc241.htm
// Also see https://en.wikipedia.org/wiki/Binomial_proportion_confidence_interval#Wilson_score_interval
func Wilson(success, fail int, z float64) float64 {
	ns := float64(success)
	nf := float64(fail)
	n := ns + nf
	z2 := z * z
	return (ns+z2/2)/(n+z2) + (z/(n+z2))*math.Sqrt((ns*nf)/n+z2/4)
}

func IsTooEasy(correct, incorrect int) bool {
	// Threshold can't be too high or the tuner will be too conservative.
	// Only uses 0.80 confidence, higher values require too many samples.

	z := -0.845 // z-score for one-sided confidence interval (80% confidence)
	lower := Wilson(correct, incorrect, z)

	// 80% likelihood that the true proportion is bounded below by `lower`.
	// > 0.9 test is too slow when incorrect > 0
	return lower > 0.875
}

func IsTooHard(correct, incorrect int) bool {
	z := 2.325 // z-score for one-sided confidence interval
	upper := Wilson(correct, incorrect, z)

	// 99% confident that the true proportion is bounded above by `upper`
	return upper < 0.8
}
