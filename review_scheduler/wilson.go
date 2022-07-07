// For computing Wilson score intervals (used by autoTune)
package review_scheduler

import (
	"math"
)

// Computes a bounary point of a Wilson score interval.
// See https://www.itl.nist.gov/div898/handbook/prc/section2/prc241.htm
// Also see https://en.wikipedia.org/wiki/Binomial_proportion_confidence_interval#Wilson_score_interval
func Wilson(success, fail int, z float64) float64 {
	ns := float64(success)
	nf := float64(fail)
	n := ns + nf
	z2 := z * z
	return (ns+z2/2)/(n+z2) + (z/(n+z2))*math.Sqrt((ns*nf)/n+z2/4)
}
