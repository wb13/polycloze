// Item buffer

import { fetchFlashcards } from "./api";
import { PartWithAnswers, hasAnswers } from "./blank";
import { Difficulty, DifficultyTuner } from "./difficulty";
import { Item } from "./item";
import { Sentence } from "./sentence";
import { edit } from "./unsaved";

type ReviewResult = {
  word: string;
  correct: boolean;
};

function* getBlankParts(sentence: Sentence): IterableIterator<PartWithAnswers> {
  for (const part of sentence.parts) {
    if (hasAnswers(part)) {
      yield part as PartWithAnswers;
    }
  }
}

export class ItemBuffer {
  buffer: Item[];
  keys: Set<string>;
  difficultyTuner: DifficultyTuner;
  reviews: ReviewResult[];

  constructor(difficulty: Difficulty = {}) {
    this.difficultyTuner = new DifficultyTuner(difficulty);
    this.buffer = [];
    this.keys = new Set();
    this.reviews = [];

    const listener = (event: Event) => {
      const review = (event as CustomEvent).detail;
      this.reviews.push(review);

      const changed = this.difficultyTuner.update(review.correct);
      if (changed) {
        // Buffered flashcards become stale when user's estimated level
        // changes.
        // Reduce number of flashcards in the buffer to trigger a refill.
        // Leaves some items in the buffer to avoid waiting for new flashcards.
        for (const key of this.buffer.splice(3)) {
          // Find word.
          let word: string = "";
          const parts = key.sentence.parts;
          for (const part of key.sentence.parts) {
            const answers = part.answers;
            if (answers != null && answers.length > 0) {
              word = answers[0].normalized;
              break;
            }
          }
          this.keys.delete(word);
        }
      }
    };

    // NOTE this never gets removed
    window.addEventListener("polycloze-review", listener);
  }

  // Add item if it's not a duplicate.
  add(item: Item): boolean {
    const parts = Array.from(getBlankParts(item.sentence));

    const words: string[] = [];
    const isDuplicate = parts.some((part) => {
      const word = part.answers[0].normalized;
      words.push(word);
      return this.keys.has(word);
    });
    if (isDuplicate) {
      return false;
    }
    this.buffer.push(item);
    words.forEach((word) => this.keys.add(word));
    return true;
  }

  // Fetches flashcards from the server and stores them in the buffer.
  async fetch(limit: number): Promise<Item[]> {
    const save = edit();

    const reviews = this.reviews.splice(0);
    const { items } = await fetchFlashcards({
      limit,
      reviews,
      exclude: Array.from(this.keys),
    });
    items.forEach((item) => this.add(item));
    this.reviews.forEach((review) => this.keys.delete(review.word));

    save();
    return items;
  }

  // Returns Promise<Item>.
  // May return undefined if there are no items left for review and there are
  // no new items left.
  async take(): Promise<Item | undefined> {
    let promise = null;
    if (this.buffer.length < 3) {
      // Could be up to 50 (fewer requests sent to server), but let's leep it
      // at 10 to minimize lost progress in case the browser closes abruptly.
      promise = this.fetch(10);
    }
    if (this.buffer.length <= 0) {
      await promise;
    }
    return this.buffer.shift();
  }
}

// Dispatches custom event to tell item buffer about review result.
export function announceResult(result: ReviewResult) {
  const event = new CustomEvent("polycloze-review", {
    detail: result,
  });
  window.dispatchEvent(event);
}
