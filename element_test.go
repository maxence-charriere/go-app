package app

import (
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app/markup"
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
			d.test.Log(err)
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

func (w *testWindow) Load(rawurl string) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	componame, ok := markup.ComponentNameFromURL(u)
	if !ok {
		return nil
	}

	compo, err := w.compoBuilder.New(componame)
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

func (w *testWindow) Center() {
}

func (w *testWindow) Size() (width, height float64) {
	return
}

func (w *testWindow) Resize(width, height float64) {
}

func (w *testWindow) Focus() {
	w.lastFocus = time.Now()
}

func (w *testWindow) ToggleFullScreen() {
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
			d.test.Log(err)
		}
	}
	return menu
}

func (m *testMenu) ID() uuid.UUID {
	return m.id
}

func (m *testMenu) Load(rawurl string) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	componame, ok := markup.ComponentNameFromURL(u)
	if !ok {
		return nil
	}

	compo, err := m.compoBuilder.New(componame)
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

func TestElementStore(t *testing.T) {

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "should add an element",
			test: testElementStoreAdd,
		},
		{
			name: "should add an element with components",
			test: testElementStoreAddElementWithComponent,
		},
		{
			name: "should fail to add an element when full",
			test: testElementStoreAddWhenFull,
		},
		{
			name: "add element with same id should fail",
			test: testElementStoreAddElementWithSameID,
		},
		{
			name: "should remove an element",
			test: testElementStoreRemove,
		},
		{
			name: "should get an element",
			test: testElementStoreElement,
		},
		{
			name: "should not get an element",
			test: testElementStoreElementNotFound,
		},
		{
			name: "should get an element by component",
			test: testElementStoreElementByComponent,
		},
		{
			name: "should not get an element by component",
			test: testElementStoreElementByComponentNotFound,
		},
		{
			name: "should sort the elements with components",
			test: testElementStoreSort,
		},
		{
			name: "should return the number of elements",
			test: testElementStoreLen,
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testElementStoreAdd(t *testing.T) {
	store := newElementStore(42)

	if err := store.Add(&testElement{
		id: uuid.New(),
	}); err != nil {
		t.Fatal(err)
	}

	if l := len(store.elements); l != 1 {
		t.Error("store should have 1 element:", l)
	}
	if l := len(store.elementsWithComponents); l != 0 {
		t.Error("store should not have an element with components")
	}
}

func testElementStoreAddElementWithComponent(t *testing.T) {
	store := newElementStore(42)

	if err := store.Add(&testWindow{
		id: uuid.New(),
	}); err != nil {
		t.Fatal(err)
	}

	if l := len(store.elements); l != 1 {
		t.Error("store should have 1 element:", l)
	}
	if l := len(store.elementsWithComponents); l != 1 {
		t.Error("store should have 1 element with components:", l)
	}
}

func testElementStoreAddElementWithSameID(t *testing.T) {
	store := newElementStore(42)
	window := &testWindow{
		id: uuid.New(),
	}

	if err := store.Add(window); err != nil {
		t.Fatal(err)
	}

	err := store.Add(window)
	if err == nil {
		t.Fatal("should not add a window twice")
	}
	t.Log()

}

func testElementStoreAddWhenFull(t *testing.T) {
	store := newElementStore(42)

	newWindow := func() *testWindow {
		return &testWindow{
			id: uuid.New(),
		}
	}

	for i := 0; i < store.capacity; i++ {
		if err := store.Add(newWindow()); err != nil {
			t.Fatal(err)
		}
	}

	err := store.Add(newWindow())
	if err == nil {
		t.Fatal("adding an element should return an error")
	}
	t.Log(err)
}

func testElementStoreRemove(t *testing.T) {
	store := newElementStore(42)
	window := &testWindow{
		id: uuid.New(),
	}

	if err := store.Add(window); err != nil {
		t.Fatal(err)
	}

	store.Remove(window)

	if l := len(store.elements); l != 0 {
		t.Error("store should not have elements:", l)
	}
	if l := len(store.elementsWithComponents); l != 0 {
		t.Error("store should not have elements with components:", l)
	}
}

func testElementStoreElement(t *testing.T) {
	store := newElementStore(42)
	window := &testWindow{
		id: uuid.New(),
	}

	if err := store.Add(window); err != nil {
		t.Fatal(err)
	}

	elem, ok := store.Element(window.ID())
	if !ok {
		t.Fatalf("no element with id %v found", window.ID())
	}
	if elem != window {
		t.Fatal("element should be the window")
	}
}

func testElementStoreElementNotFound(t *testing.T) {
	store := newElementStore(42)
	if _, ok := store.Element(uuid.New()); ok {
		t.Fatal("no element should have been found")
	}
}

func testElementStoreElementByComponent(t *testing.T) {
	compoBuilder := markup.NewCompoBuilder()
	if err := compoBuilder.Register(&Component{}); err != nil {
		t.Fatal(err)
	}

	compo := &Component{}
	env := markup.NewEnv(compoBuilder)
	if _, err := env.Mount(compo); err != nil {
		t.Fatal(err)
	}

	window := &testWindow{
		id:           uuid.New(),
		compoBuilder: compoBuilder,
		env:          env,
	}

	store := newElementStore(42)
	if err := store.Add(window); err != nil {
		t.Fatal(err)
	}

	elem, err := store.ElementByComponent(compo)
	if err != nil {
		t.Fatal(err)
	}
	if elem != window {
		t.Fatal("element should be the window")
	}
}

func testElementStoreElementByComponentNotFound(t *testing.T) {
	store := newElementStore(42)

	if _, err := store.ElementByComponent(&Component{}); err == nil {
		t.Fatal("no element should have been found")
	}
}

func testElementStoreSort(t *testing.T) {
	store := newElementStore(42)

	for i := 0; i < 10; i++ {
		if err := store.Add(&testMenu{
			id:        uuid.New(),
			lastFocus: time.Now(),
		}); err != nil {
			t.Fatal(err)
		}
	}

	window := &testWindow{
		id:        uuid.New(),
		lastFocus: time.Now(),
	}
	if err := store.Add(window); err != nil {
		t.Fatal(err)
	}

	elements := store.elementsWithComponents
	for i, elem := range elements {
		if elem.ID() == window.ID() {
			elements[i], elements[5] = elements[5], elements[i]
			break
		}
	}

	store.Sort()

	if elem := store.elementsWithComponents[0]; elem != window {
		t.Fatalf("1st element with components should be the window: %T", elem)
	}
}

func testElementStoreLen(t *testing.T) {
	store := newElementStore(42)

	for i := 0; i < 10; i++ {
		if err := store.Add(&testMenu{
			id:        uuid.New(),
			lastFocus: time.Now(),
		}); err != nil {
			t.Fatal(err)
		}
	}

	if l := store.Len(); l != 10 {
		t.Fatal("store should have 10 elements:", l)
	}
}
