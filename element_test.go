package app

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app/markup"
)

type element struct {
	id uuid.UUID
}

func newElement() *element {
	return &element{
		id: uuid.New(),
	}
}

func (e *element) ID() uuid.UUID {
	return e.id
}

type elementWithComponent struct {
	id        uuid.UUID
	lastFocus time.Time
}

func newElementWithComponent() *elementWithComponent {
	return &elementWithComponent{
		id:        uuid.New(),
		lastFocus: time.Now(),
	}
}

func (e *elementWithComponent) ID() uuid.UUID {
	return e.id
}

func (e *elementWithComponent) Load(url string) error {
	return nil
}

func (e *elementWithComponent) Contains(c markup.Component) bool {
	return false
}

func (e *elementWithComponent) Render(c markup.Component) error {
	return nil
}

func (e *elementWithComponent) LastFocus() time.Time {
	return e.lastFocus
}

func TestElementStoreAdd(t *testing.T) {
	capacity := 42
	store := newElementStore(capacity)
	var lastElem ElementWithComponent

	for i := 0; i < capacity; i++ {
		lastElem = newElementWithComponent()
		if err := store.Add(lastElem); err != nil {
			t.Fatal(err)
		}
	}

	if firstElem := store.elementsWithComponents[0]; firstElem != lastElem {
		t.Fatal("last element should have moved to be the first element")
	}

	err := store.Add(newElement())
	if err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)
}

func TestElementStoreDelete(t *testing.T) {
	capacity := 42
	store := newElementStore(capacity)

	elem := newElement()
	if err := store.Add(elem); err != nil {
		t.Fatal(err)
	}
	store.Remove(elem)

	elemWithCompo := newElementWithComponent()
	if err := store.Add(elemWithCompo); err != nil {
		t.Fatal(err)
	}
	store.Remove(elemWithCompo)

	if len(store.elements) != 0 {
		t.Error("store.elements should be empty")
	}
	if len(store.elementsWithComponents) != 0 {
		t.Error("store.elementsWithComponents should be empty")
	}
}
