package core

import (
	"sync"
)

// History represents a store that contains URL ordered chronologically.
type History struct {
	once    sync.Once
	mutex   sync.RWMutex
	index   int
	history []string
}

func (h *History) init() {
	h.index = -1
	h.history = make([]string, 0, 32)
}

// Len returns the number of entries recorded in the history.
func (h *History) Len() int {
	h.once.Do(h.init)
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return len(h.history)
}

// Current returns the current entry.
func (h *History) Current() (url string) {
	h.once.Do(h.init)

	if h.Len() == 0 {
		return ""
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	url = h.history[h.index]
	return url
}

// NewEntry adds an entry to the history.
// The entry is added after the current one.
// All the entries that was after the current one are removed.
func (h *History) NewEntry(url string) {
	h.once.Do(h.init)

	if len(url) == 0 {
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	var history []string

	if len(h.history) == 0 {
		history = h.history
	} else {
		history = h.history[:h.index+1]
	}

	h.history = append(history, url)
	h.index++
}

// CanPrevious reports whether there is a previous entry.
func (h *History) CanPrevious() bool {
	h.once.Do(h.init)
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.index > 0
}

// Previous returns the previous entry.
func (h *History) Previous() (url string) {
	h.once.Do(h.init)

	if !h.CanPrevious() {
		return ""
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.index--
	url = h.history[h.index]
	return url
}

// CanNext reports whether there is a next entry.
func (h *History) CanNext() bool {
	h.once.Do(h.init)
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.index < h.Len()-1
}

// Next returns the next entry.
func (h *History) Next() (url string) {
	h.once.Do(h.init)

	if !h.CanNext() {
		return ""
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.index++
	url = h.history[h.index]
	return url
}
