// Difficulty tuner.

import { isTooEasy, isTooHard } from "./wilson";

export class DifficultyTuner {
  level: number;
  correct: number;
  incorrect: number;
  min: number;
  max: number;

  constructor(level = 0, correct = 0, incorrect = 0, min = 0, max = 0) {
    this.level = level;
    this.correct = correct;
    this.incorrect = incorrect;
    this.min = min;
    this.max = max;
  }

  // Updates level statistics.
  // Returns true if level changed.
  // Also resets `correct` and `incorrect` counters if so.
  update(correct: boolean): boolean {
    if (correct) {
      this.correct++;
    } else {
      this.incorrect++;
    }

    const level = this.level;
    if (isTooEasy(this.correct, this.incorrect)) {
      this.level = Math.min(level + 1, this.max);
      if (level === this.level) {
        return false;
      }
      this.correct = 0;
      this.incorrect = 0;
      return true;
    }

    if (isTooHard(this.correct, this.incorrect)) {
      this.level = Math.max(level - 1, this.min);
      if (level === this.level) {
        return false;
      }
      this.correct = 0;
      this.incorrect = 0;
      return true;
    }

    return false;
  }
}
