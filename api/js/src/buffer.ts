// Item buffer

import { fetchFlashcards, sendReviewResults, fetchSentences } from "./api";
import { PartWithAnswers, hasAnswers } from "./blank";
import { Difficulty, DifficultyTuner } from "./difficulty";
import { Item } from "./item";
import { ReviewResult, RandomSentence } from "./schema";
import { Sentence } from "./sentence";

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

      if (!review.new) {
        // Only tune difficulty based on new words.
        return;
      }

      // TODO new word might have frequencyClass < current level
      // But does it matter if all difficult words have already been seen?
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

    // MDN recommends `visibilitychange` instead of `unload` and `beforeunload`
    // because `visibilitychange` is more reliable on mobile.
    window.addEventListener("visibilitychange", () => {
      if (document.visibilityState === "hidden" && this.reviews.length > 0) {
        sendReviewResults(this.reviews, this.difficultyTuner.difficulty);

        // Clear buffer to avoid resending sent reviews in case the user
        // switches back to this tab.
        const reviews = this.reviews.splice(0);
        reviews.forEach((review) => this.keys.delete(review.word));

        // This doesn't update the difficulty stats, but it shouldn't be a
        // problem because the item buffer will get the updated stats on the
        // next fetch.
      }
    });
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
    const reviews = this.reviews.splice(0);
    const { items, difficulty } = await fetchFlashcards({
      limit,
      reviews,
      difficulty: this.difficultyTuner.difficulty,
      exclude: Array.from(this.keys),
    });
    items.forEach((item) => this.add(item));
    reviews.forEach((review) => this.keys.delete(review.word));
    this.difficultyTuner.reset(difficulty);
    return items;
  }

  // Returns Promise<Item>.
  // May return undefined if there are no items left for review and there are
  // no new items left.
  async take(): Promise<Item | undefined> {
    let promise = null;
    if (this.buffer.length < 3) {
      // How many flashcards to fetch?
      // The fewer flashcards you fetch, the fewer the flashcards that get
      // wasted when the student's estimated level changes.
      // But if you fetch more flashcards per batch, you have to contact the
      // server less frequently...
      promise = this.fetch(30);
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

export class RandomSentenceBuffer {
  buffer: RandomSentence[];
  difficulty: number;

  constructor(difficulty: number) {
    this.buffer = [];
    this.difficulty = difficulty;
  };

  async fetch(limit: number): Promise<Item[]> {
    const items = await fetchSentences({difficulty: this.difficulty, limit: limit});
    items.forEach((item) => this.buffer.push(item));
    return items;
  }

  async take(): Promise<RandomSentence | undefined> {
    let promise = null;
    if (this.buffer.length < 3) {
      promise = await this.fetch(30);
    }
    if (this.buffer.length <= 0) {
      await promise;
    }
    return this.buffer.shift();
  }
}

