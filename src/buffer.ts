// Item buffer

import { fetchItems } from "./data";
import { Item } from "./item";
import { Sentence } from "./sentence";

function * oddParts(sentence: Sentence): IterableIterator<string> {
    for (const [i, part] of sentence.parts.entries()) {
        if (i % 2 === 1) {
            yield part;
        }
    }
}

export class ItemBuffer {
    buffer: Item[];
    keys: Set<string>;
    backgroundFetch: Promise<Item[]> | null;

    constructor() {
        this.buffer = [];
        this.keys = new Set();
        this.backgroundFetch = null;

        const listener = (event: Event) => {
            const word = (event as CustomEvent).detail.word;
            this.keys.delete(word);
        };

        // NOTE this never gets removed
        window.addEventListener("polycloze-unbuffer", listener);
    }

    // Add item if it's not a duplicate.
    add(item: Item): boolean {
        const parts = Array.from(oddParts(item.sentence));
        if (parts.some(part => this.keys.has(part))) {
            return false;
        }
        this.buffer.push(item);
        parts.forEach(part => this.keys.add(part));
        return true;
    }

    // Returns Promise Item and a function should be called after submitReview.
    async take(): Promise<Item> {
        if (this.backgroundFetch != null) {
            const items = await this.backgroundFetch;
            this.backgroundFetch = null;
            items.forEach(item => this.add(item));
        }

        if (this.buffer.length === 0) {
            this.backgroundFetch = fetchItems(2, Array.from(this.keys));
        } else if (this.buffer.length < 20) {
            this.backgroundFetch = fetchItems(10, Array.from(this.keys));
        }

        if (this.buffer.length === 0) {
            return this.take();
        }
        return this.buffer.shift()!;
    }
}

export function dispatchUnbuffer(word: string) {
    const event = new CustomEvent("polycloze-unbuffer", {
        detail: { word }
    });
    window.dispatchEvent(event);
}
