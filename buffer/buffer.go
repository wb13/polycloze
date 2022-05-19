// Buffered item generator.
package buffer

import (
	"sync"

	"github.com/lggruspe/polycloze/flashcards"
)

type ItemBuffer struct {
	Channel chan flashcards.Item
	words   map[string]bool
	mutex   sync.Mutex // for words map

	ig flashcards.ItemGenerator
}

func NewItemBuffer(ig flashcards.ItemGenerator, capacity int) ItemBuffer {
	return ItemBuffer{
		Channel: make(chan flashcards.Item, capacity),
		words:   make(map[string]bool),
		ig:      ig,
	}
}

func (buf *ItemBuffer) Add(word string) {
	buf.mutex.Lock()
	buf.words[word] = true
	buf.mutex.Unlock()
}

// Fetch items and store in buffer.
func (buf *ItemBuffer) Fetch() error {
	n := cap(buf.Channel) - len(buf.Channel)
	words, err := buf.ig.GenerateWordsWith(n, func(word string) bool {
		buf.mutex.Lock()
		_, ok := buf.words[word]
		buf.mutex.Unlock()
		return !ok
	})
	if err != nil {
		return err
	}

	for _, word := range words {
		buf.Add(word)
	}
	buf.ig.GenerateItemsIntoChannel(buf.Channel, words)
	return nil
}

// Take an item out of buffer.
func (buf *ItemBuffer) Take() flashcards.Item {
	item := <-buf.Channel
	word := item.Sentence.Parts[1]
	buf.mutex.Lock()
	delete(buf.words, word)
	buf.mutex.Unlock()
	return item
}
