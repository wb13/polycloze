// Buffered item generator.
package buffer

import (
	"strings"
	"sync"

	"github.com/lggruspe/polycloze/flashcards"
)

type ItemBuffer struct {
	Channel chan flashcards.Item
	words   map[string]bool
	mutex		sync.RWMutex // for words map

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
	buf.words[strings.ToLower(word)] = true
	buf.mutex.Unlock()
}

func (buf *ItemBuffer) Contains(word string) bool {
	buf.mutex.RLock()
	_, ok := buf.words[strings.ToLower(word)]
	buf.mutex.RUnlock()
	return ok
}

// Fetch n items and store in buffer.
// Pass non-positive number if you want to fill the buffer.
func (buf *ItemBuffer) Fetch(n int) error {
	max := cap(buf.Channel) - len(buf.Channel)
	if n <= 0 || max < n {
		n = max
	}

	words, err := buf.ig.GenerateWordsWith(n, func(word string) bool {
		return !buf.Contains(word)
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

func (buf *ItemBuffer) Refill() {
	if len(buf.Channel) == 0 {
		buf.Fetch(2)
		// Quickly generate 2 flashcards.
		// Can't be Fetch(1), because worst case scenario: it keeps fetching 1
		// repeatedly.

		go buf.Fetch(20)
		// 20 flashcards take about ~2 seconds to generate, the same amount of time
		// it takes to answer a flashcard.
	}
	if 3*len(buf.Channel) <= 2*cap(buf.Channel) {
		go buf.Fetch(-1)
	}
}

// Take an item out of buffer.
func (buf *ItemBuffer) Take() flashcards.Item {
	go buf.Refill()

	item := <-buf.Channel
	word := item.Sentence.Parts[1]
	buf.mutex.Lock()
	delete(buf.words, strings.ToLower(word))
	buf.mutex.Unlock()
	return item
}

// Take many items out of buffer, where many = cap(buffer channel) / 3.
func (buf *ItemBuffer) TakeMany() []flashcards.Item {
	var items []flashcards.Item
	for i := 0; i < cap(buf.Channel)/3; i++ {
		items = append(items, buf.Take())
	}
	return items
}
