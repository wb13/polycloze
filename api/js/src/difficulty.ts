// Difficulty tuner.

import { isTooEasy, isTooHard } from "./wilson";

export type Difficulty = {
  level?: number;
  correct?: number;
  incorrect?: number;
  min?: number;
  max?: number;
};

export class DifficultyTuner {
  // @ts-ignore because typescript can't see initializer in `this.reset`.
  level: number;
  // @ts-ignore
  correct: number;
  // @ts-ignore
  incorrect: number;
  // @ts-ignore
  min: number;
  // @ts-ignore
  max: number;

  constructor(difficulty: Difficulty = {}) {
    this.reset(difficulty);
  }

  reset(difficulty: Difficulty) {
    let { level, correct, incorrect, min, max } = difficulty;
    if (min == null || min < 0) {
      min = 0;
    }
    if (max == null || max < 0) {
      max = Infinity;
    }
    if (correct == null || correct < 0) {
      correct = 0;
    }
    if (incorrect == null || incorrect < 0) {
      incorrect = 0;
    }
    if (level == null || level < 0) {
      level = min;
    }

    this.min = min;
    this.max = max;
    this.level = level;
    this.correct = correct;
    this.incorrect = incorrect;
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
