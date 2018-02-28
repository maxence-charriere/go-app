package tests

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

type element struct {
	id uuid.UUID
}

func newElement() *element {
	elem := &element{
		id: uuid.New(),
	}
	return elem
}

func (e *element) ID() uuid.UUID {
	return e.id
}

type elementWithComponent struct {
	id        uuid.UUID
	factory   app.Factory
	lastFocus time.Time
	component app.Component
}

func newElementWithComponent() *elementWithComponent {
	factory := app.NewFactory()
	factory.Register(&Foo{})

	return &elementWithComponent{
		id:        uuid.New(),
		factory:   factory,
		lastFocus: time.Now(),
	}
}

func (e *elementWithComponent) ID() uuid.UUID {
	return e.id
}

func (e *elementWithComponent) Load(rawurl string, v ...interface{}) error {
	rawurl = fmt.Sprintf(rawurl, v...)

	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	compo, err := e.factory.New(app.ComponentNameFromURL(u))
	if err != nil {
		return err
	}

	e.component = compo
	return nil
}

func (e *elementWithComponent) Component() app.Component {
	return e.component
}

func (e *elementWithComponent) Contains(c app.Component) bool {
	return e.component == c
}

func (e *elementWithComponent) Render(c app.Component) error {
	e.component = c
	return nil
}

func (e *elementWithComponent) LastFocus() time.Time {
	return e.lastFocus
}

// TestElemDB is a test suite used to ensure that all element databases
// implementations behave the same.
func TestElemDB(t *testing.T, newElemDB func() app.ElemDB) {
	tests := []struct {
		scenario string
		function func(t *testing.T, db app.ElemDB)
	}{
		{
			scenario: "adds an element",
			function: testElemDBAdd,
		},
		{
			scenario: "adds an element with components",
			function: testElemDBAddElementWithComponent,
		},
		{
			scenario: "adding element with same id returns an error",
			function: testElemDBAddElementWithSameID,
		},
		{
			scenario: "removes an element",
			function: testElemDBRemove,
		},
		{
			scenario: "get an element",
			function: testElemDBElement,
		},
		{
			scenario: "get a nonexistent element returns false",
			function: testElemDBElementNotFound,
		},
		{
			scenario: "get an element by component",
			function: testElemDBElementByComponent,
		},
		{
			scenario: "get an element by not mounted component returns an error",
			function: testElemDBElementByComponentNotFound,
		},
		{
			scenario: "sorts the elements with components",
			function: testElemDBSort,
		},
		{
			scenario: "get the number of elements",
			function: testElemDBLen,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t, newElemDB())
		})

	}
}

func testElemDBAdd(t *testing.T, db app.ElemDB) {
	if err := db.Add(newElement()); err != nil {
		t.Fatal(err)
	}
	if l := db.Len(); l != 1 {
		t.Error("database doesn't have 1 element:", l)
	}
	if len(db.ElementsWithComponents()) != 0 {
		t.Error("database have an element with component")
	}
}

func testElemDBAddElementWithComponent(t *testing.T, db app.ElemDB) {
	if err := db.Add(newElementWithComponent()); err != nil {
		t.Fatal(err)
	}
	if l := db.Len(); l != 1 {
		t.Error("database doesn't have 1 element:", l)
	}
	if l := len(db.ElementsWithComponents()); l != 1 {
		t.Error("database doesn't have 1 element:", l)
	}
}

func testElemDBAddElementWithSameID(t *testing.T, db app.ElemDB) {
	elem := newElementWithComponent()

	if err := db.Add(elem); err != nil {
		t.Fatal(err)
	}

	err := db.Add(elem)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testElemDBRemove(t *testing.T, db app.ElemDB) {
	elem := newElementWithComponent()

	if err := db.Add(elem); err != nil {
		t.Fatal(err)
	}

	db.Remove(elem)

	if l := db.Len(); l != 0 {
		t.Error("database has elements:", l)
	}
	if l := len(db.ElementsWithComponents()); l != 0 {
		t.Error("database has elements with components:", l)
	}
}

func testElemDBElement(t *testing.T, db app.ElemDB) {
	elem := newElementWithComponent()

	if err := db.Add(elem); err != nil {
		t.Fatal(err)
	}

	ret, err := db.Element(elem.ID())
	if err != nil {
		t.Fatal(err)
	}
	if ret != elem {
		t.Fatal("returned element is no the added element")
	}
}

func testElemDBElementNotFound(t *testing.T, db app.ElemDB) {
	_, err := db.Element(uuid.New())
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testElemDBElementByComponent(t *testing.T, db app.ElemDB) {
	compo := &Bar{}
	elem := newElementWithComponent()
	elem.component = compo

	if err := db.Add(elem); err != nil {
		t.Fatal(err)
	}

	ret, err := db.ElementByComponent(compo)
	if err != nil {
		t.Fatal(err)
	}
	if ret != elem {
		t.Fatal("returned element is not the added element")
	}
}

func testElemDBElementByComponentNotFound(t *testing.T, db app.ElemDB) {
	if _, err := db.ElementByComponent(&Foo{}); err == nil {
		t.Fatal("an element is found")
	}
}

func testElemDBSort(t *testing.T, db app.ElemDB) {
	for i := 0; i < 10; i++ {
		if err := db.Add(newElementWithComponent()); err != nil {
			t.Fatal(err)
		}
	}

	lastElem := newElementWithComponent()
	if err := db.Add(lastElem); err != nil {
		t.Fatal(err)
	}

	db.Sort()

	if elem := db.ElementsWithComponents()[0]; elem != lastElem {
		t.Fatalf("1st element with components is not the last added element: %T", elem)
	}
}

func testElemDBLen(t *testing.T, db app.ElemDB) {
	for i := 0; i < 10; i++ {
		if err := db.Add(newElementWithComponent()); err != nil {
			t.Fatal(err)
		}
	}

	if l := db.Len(); l != 10 {
		t.Fatal("elemDB doesn't have 10 elements:", l)
	}
}
