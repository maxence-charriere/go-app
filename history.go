package app

import (
	"sync"

	"github.com/pkg/errors"
)

// History is the interface that describes a store that contains URL ordered
// chronologically.
type History interface {
	// Len returns the number of entries recorded in the history.
	Len() int

	// Current returns the current entry.
	Current() (url string, err error)

	// NewEntry adds an entry to the history.
	// The entry is added after the current one.
	// All the entries that was after the current one are removed.
	NewEntry(url string)

	// CanPrevious reports whether there is a previous entry.
	CanPrevious() bool

	// Previous returns the previous entry.
	Previous() (url string, err error)

	// CanNext reports whether there is a next entry.
	CanNext() bool

	// Next returns the next entry.
	Next() (url string, err error)
}

// NewHistory creates an history.
func NewHistory() History {
	return &history{
		index:   -1,
		history: make([]string, 0, 32),
	}
}

type history struct {
	index   int
	history []string
}

func (h *history) Len() int {
	return len(h.history)
}

func (h *history) Current() (url string, err error) {
	if h.Len() == 0 {
		return "", errors.New("history does not have entries")
	}

	url = h.history[h.index]
	return url, nil
}

func (h *history) NewEntry(url string) {
	var history []string

	if h.Len() == 0 {
		history = h.history
	} else {
		history = h.history[:h.index+1]
	}

	h.history = append(history, url)
	h.index++
}

func (h *history) CanPrevious() bool {
	return h.index > 0
}

func (h *history) Previous() (url string, err error) {
	if !h.CanPrevious() {
		return "", errors.New("history does not have a previous entry to return")
	}

	h.index--
	url = h.history[h.index]
	return url, nil
}

func (h *history) CanNext() bool {
	return h.index < h.Len()-1
}

func (h *history) Next() (url string, err error) {
	if !h.CanNext() {
		return "", errors.New("history does not have a next entry to return")
	}

	h.index++
	url = h.history[h.index]
	return url, nil
}

// ConcurrentHistory returns a decorated version of the given history that
// is safe for concurrent operations.
func ConcurrentHistory(history History) History {
	return &concurrentHistory{
		base: history,
	}
}

type concurrentHistory struct {
	base  History
	mutex sync.Mutex
}

func (h *concurrentHistory) Len() int {
	h.mutex.Lock()
	l := h.base.Len()
	h.mutex.Unlock()
	return l
}

func (h *concurrentHistory) Current() (url string, err error) {
	h.mutex.Lock()
	url, err = h.base.Current()
	h.mutex.Unlock()
	return url, err
}

func (h *concurrentHistory) NewEntry(url string) {
	h.mutex.Lock()
	h.base.NewEntry(url)
	h.mutex.Unlock()
}

func (h *concurrentHistory) CanPrevious() bool {
	h.mutex.Lock()
	ok := h.base.CanPrevious()
	h.mutex.Unlock()
	return ok
}

func (h *concurrentHistory) Previous() (url string, err error) {
	h.mutex.Lock()
	url, err = h.base.Previous()
	h.mutex.Unlock()
	return url, err
}

func (h *concurrentHistory) CanNext() bool {
	h.mutex.Lock()
	ok := h.base.CanNext()
	h.mutex.Unlock()
	return ok
}

func (h *concurrentHistory) Next() (url string, err error) {
	h.mutex.Lock()
	url, err = h.base.Next()
	h.mutex.Unlock()
	return url, err
}
