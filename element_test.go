package app

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app/markup"
	"github.com/murlokswarm/app/url"
	"github.com/pkg/errors"
)

// A simple element implementation for tests.
type testElement struct {
	id uuid.UUID
}

func newTestElement(d *testDriver) *testElement {
	elem := &testElement{
		id: uuid.New(),
	}
	d.elements.Add(elem)
	return elem
}

func (e *testElement) ID() uuid.UUID {
	return e.id
}

// A window implementation for tests.
type testWindow struct {
	config       WindowConfig
	id           uuid.UUID
	compoBuilder markup.CompoBuilder
	env          markup.Env
	lastFocus    time.Time

	onLoad  func(c markup.Component)
	onClose func()
}

func newTestWindow(d *testDriver, c WindowConfig) *testWindow {
	window := &testWindow{
		config:       c,
		id:           uuid.New(),
		compoBuilder: d.compoBuilder,
		env:          markup.NewEnv(d.compoBuilder),
		lastFocus:    time.Now(),
	}

	d.elements.Add(window)
	window.onClose = func() {
		d.elements.Remove(window)
	}

	if d.onWindowLoad != nil {
		window.onLoad = func(c markup.Component) {
			d.onWindowLoad(window, c)
		}
	}

	if len(c.DefaultURL) != 0 {
		if err := window.Load(c.DefaultURL); err != nil {
			d.Test.Log(errors.Wrap(err, ""))
		}
	}
	return window
}

func (w *testWindow) Close() {
	w.onClose()
}

func (w *testWindow) ID() uuid.UUID {
	return w.id
}

func (w *testWindow) Contains(c markup.Component) bool {
	return w.env.Contains(c)
}

func (w *testWindow) Render(c markup.Component) error {
	_, err := w.env.Update(c)
	return err
}

func (w *testWindow) LastFocus() time.Time {
	return w.lastFocus
}

func (w *testWindow) Load(u string) error {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}

	compoName, ok := parsedURL.Component()
	if !ok {
		return nil
	}

	compo, err := w.compoBuilder.New(compoName)
	if err != nil {
		return err
	}

	if _, err = w.env.Mount(compo); err != nil {
		return errors.Wrapf(err, "loading %s in test window %p failed", u, w)
	}

	if w.onLoad != nil {
		w.onLoad(compo)
	}
	return nil
}

func (w *testWindow) CanPrevious() bool {
	return false
}

func (w *testWindow) Previous() error {
	return nil
}

func (w *testWindow) CanNext() bool {
	return false
}

func (w *testWindow) Next() error {
	return nil
}

func (w *testWindow) Position() (x, y float64) {
	return
}

func (w *testWindow) Move(x, y float64) {
}

func (w *testWindow) Size() (width, height float64) {
	return
}

func (w *testWindow) Resize(width, height float64) {
}

func (w *testWindow) Focus() {
	w.lastFocus = time.Now()
}

// A menu implementation for tests.
type testMenu struct {
	config       MenuConfig
	id           uuid.UUID
	compoBuilder markup.CompoBuilder
	env          markup.Env
	lastFocus    time.Time
}

func newTestMenu(d *testDriver, c MenuConfig) *testMenu {
	menu := &testMenu{
		id:           uuid.New(),
		compoBuilder: d.compoBuilder,
		env:          markup.NewEnv(d.compoBuilder),
		lastFocus:    time.Now(),
	}
	d.elements.Add(menu)

	if len(c.DefaultURL) != 0 {
		if err := menu.Load(c.DefaultURL); err != nil {
			d.Test.Log(errors.Wrap(err, ""))
		}
	}
	return menu
}

func (m *testMenu) ID() uuid.UUID {
	return m.id
}

func (m *testMenu) Load(u string) error {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}

	compoName, ok := parsedURL.Component()
	if !ok {
		return nil
	}

	compo, err := m.compoBuilder.New(compoName)
	if err != nil {
		return err
	}

	if _, err = m.env.Mount(compo); err != nil {
		return errors.Wrapf(err, "loading %s in test menu %p failed", u, m)
	}
	return nil
}

func (m *testMenu) Contains(c markup.Component) bool {
	return m.env.Contains(c)
}

func (m *testMenu) Render(c markup.Component) error {
	_, err := m.env.Update(c)
	return err
}

func (m *testMenu) LastFocus() time.Time {
	return m.lastFocus
}

// A dock tile implementation for tests.
type testDockTile struct {
	testMenu
}

func newDockTile(d *testDriver) *testDockTile {
	dock := &testDockTile{
		testMenu: testMenu{
			id:           uuid.New(),
			compoBuilder: d.compoBuilder,
			env:          markup.NewEnv(d.compoBuilder),
			lastFocus:    time.Now(),
		},
	}
	d.elements.Add(dock)
	return dock
}

func (d *testDockTile) SetIcon(name string) error {
	return nil
}

func (d *testDockTile) SetBadge(v interface{}) {
}

func TestElementStoreAdd(t *testing.T) {
	capacity := 42
	store := newElementStore(capacity)
	var lastElem ElementWithComponent

	for i := 0; i < capacity; i++ {
		lastElem = &testWindow{
			id:        uuid.New(),
			lastFocus: time.Now(),
		}
		if err := store.Add(lastElem); err != nil {
			t.Fatal(err)
		}
	}

	if firstElem := store.elementsWithComponents[0]; firstElem != lastElem {
		t.Fatal("last element should have moved to be the first element")
	}

	overElem := &testElement{
		id: uuid.New(),
	}
	err := store.Add(overElem)
	if err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)
}

func TestElementStoreDelete(t *testing.T) {
	capacity := 42
	store := newElementStore(capacity)

	elem := &testElement{
		id: uuid.New(),
	}
	if err := store.Add(elem); err != nil {
		t.Fatal(err)
	}
	store.Remove(elem)

	elemWithCompo := &testMenu{
		id:        uuid.New(),
		lastFocus: time.Now(),
	}
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
