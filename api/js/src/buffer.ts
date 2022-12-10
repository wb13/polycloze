// Item buffer

import { fetchItems } from "./api";
import { Part } from "./blank";
import { Item } from "./item";
import { Sentence } from "./sentence";

function * oddParts(sentence: Sentence): IterableIterator<Part> {
    for (const [i, part] of sentence.parts.entries()) {
        if (i % 2 === 1) {
            yield part;
        }
    }
}

export class ItemBuffer {
    buffer: Item[];
    keys: Set<string>;

    frequencyClass?: number;

    constructor() {
        this.buffer = [];
        this.keys = new Set();

        const listener = (event: Event) => {
            const word = (event as CustomEvent).detail.word;
            // TODO make sure something gets deleted
            this.keys.delete(word);
        };

        // NOTE this never gets removed
        window.addEventListener("polycloze-unbuffer", listener);
    }

    // Add item if it's not a duplicate.
    add(item: Item): boolean {
        const parts = Array.from(oddParts(item.sentence));

        const words: string[] = [];
        const isDuplicate = parts.some(part => {
            const answers = part.answers;
            if (answers == null || answers.length === 0) {
                // Ignore part if not a blank.
                // TODO don't do it like this
                return true;
            }
            const word = answers[0].normalized;
            words.push(word);
            return this.keys.has(word);
        });
        if (isDuplicate) {
            return false;
        }
        this.buffer.push(item);
        words.forEach(word => this.keys.add(word));
        return true;
    }

    backgroundFetch(count: number) {
        setTimeout(async() => {
            const items = await fetchItems({
                n: count,
                x: Array.from(this.keys),
            });
            items.forEach(item => this.add(item));
        });
    }

    // Returns Promise<Item>.
    // May return undefined if there are no items left for review and there are
    // no new items left.
    async take(): Promise<Item | undefined> {
        if (this.buffer.length === 0) {
            const items = await fetchItems({
                n: 2,
                x: Array.from(this.keys),
            });
            this.backgroundFetch(2);
            items.forEach(item => this.add(item));
            return this.buffer.shift();
        }
        if (this.buffer.length < 10) {
            this.backgroundFetch(10);
        }
        return this.buffer.shift();
    }

    clearIfStale(frequencyClass: number) {
        if (this.frequencyClass != undefined && this.frequencyClass != frequencyClass) {
            // Leaves some items in the buffer so flashcards come continuously.
            this.buffer.splice(3);
        }
        this.frequencyClass = frequencyClass;
    }
}

export function dispatchUnbuffer(word: string) {
    const event = new CustomEvent("polycloze-unbuffer", {
        detail: { word }
    });
    window.dispatchEvent(event);
}
