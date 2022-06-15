// Item buffer

import { fetchItems } from './data'
import { Item } from './item'
import { Sentence } from './sentence'

function * oddParts (sentence: Sentence): IterableIterator<string> {
  for (const [i, part] of sentence.parts.entries()) {
    if (i % 2 === 1) {
      yield part
    }
  }
}

export class ItemBuffer {
  buffer: Item[]
  keys: Set<string>
  backgroundFetch: Promise<Item[]>

  constructor () {
    this.buffer = []
    this.keys = new Set()
    this.backgroundFetch = null
  }

  // Add item if it's not a duplicate.
  add (item: Item): boolean {
    // TODO not perfect, because no case-folding
    const parts = Array.from(oddParts(item.sentence))
    if (parts.some(part => this.keys.has(part))) {
      return false
    }
    this.buffer.push(item)
    parts.forEach(part => this.keys.add(part))
    return true
  }

  deleteParts (item: Item) {
    for (const part of oddParts(item.sentence)) {
      this.keys.delete(part)
    }
  }

  // Returns Promise Item and a function should be called after submitReview.
  async take (): Promise<[Item, () => void]> {
    if (this.backgroundFetch != null) {
      const items = await this.backgroundFetch
      this.backgroundFetch = null
      items.forEach(item => this.add(item))
    }

    if (this.buffer.length < 20) {
      this.backgroundFetch = fetchItems(10)
    }
    if (this.buffer.length === 0) {
      return this.take()
    }

    const item = this.buffer.shift()
    return [item, () => this.deleteParts(item)]
  }
}
