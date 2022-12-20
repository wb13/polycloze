// For computing Wilson score intervals.

// Computes a boundary point of a Wilson score interval.
// See:
// - polycloze/wilson/wilson.go
// - https://www.itl.nist.gov/div898/handbook/prc/section2/prc241.htm
// - https://en.wikipedia.org/wiki/Binomial_proportion_confidence_interval#Wilson_score_interval
function wilson(success: number, fail: number, z: number): number {
  const [ns, nf, n] = [success, fail, success + fail];
  const z2 = z ** 2;
  return (
    (ns + z2 / 2) / (n + z2) +
    (z / (n + z2)) * Math.sqrt((ns * nf) / n + z2 / 4)
  );
}

export function isTooEasy(correct: number, incorrect: number): boolean {
  // Threshold can't be too high or the tuner will be too conservative.
  // Only uses 0.85 confidence, higher values require too many samples.

  // z-score for one-sided confidence interval (85% confidence)
  const z = -1.035;
  const lower = wilson(correct, incorrect, z);

  // 85% likelihood that the true proportion is bounded below by `lower`.
  // It's too hard to level up with a 0.9 test when incorrect > 0.
  return lower > 0.85;

  // 0.85 threshold is chosen so tuner won't trigger with < 5 samples.
}

export function isTooHard(correct: number, incorrect: number): boolean {
  // z-score for one-sided confidence interval.
  const z = 3.1;
  const upper = wilson(correct, incorrect, z);

  // 99.9% confident that the true proportion is bounded above by `upper`.
  return upper < 0.8;
}
