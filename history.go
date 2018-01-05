package app

import "github.com/pkg/errors"

// History is a store that contains URL ordered chronologically.
type History struct {
	index   int
	history []string
}

// NewHistory creates an history.
func NewHistory() *History {
	return &History{
		index:   -1,
		history: make([]string, 0, 32),
	}
}

// Len returns the number of entries recorded in the history.
func (h *History) Len() int {
	return len(h.history)
}

// Current returns the current entry.
func (h *History) Current() (url string, err error) {
	if h.Len() == 0 {
		err = errors.New("no entry")
		return
	}

	url = h.history[h.index]
	return
}

// NewEntry adds an entry to the history.
// Then entry is added after the current one.
// All the entries that was after the current one are removed.
func (h *History) NewEntry(url string) {
	var history []string

	if h.Len() == 0 {
		history = h.history
	} else {
		history = h.history[:h.index+1]
	}

	h.history = append(history, url)
	h.index++
}

// CanPrevious reports whether there is a previous entry.
func (h *History) CanPrevious() bool {
	return h.index > 0
}

// Previous returns the previous entry.
func (h *History) Previous() (url string, err error) {
	if !h.CanPrevious() {
		err = errors.New("no entry to go back")
		return
	}

	h.index--
	url = h.history[h.index]
	return
}

// CanNext reports whether there is a next entry.
func (h *History) CanNext() bool {
	return h.index < h.Len()-1
}

// Next returns the next entry.
func (h *History) Next() (url string, err error) {
	if !h.CanNext() {
		err = errors.New("no entry to go next")
		return
	}

	h.index++
	url = h.history[h.index]
	return
}
