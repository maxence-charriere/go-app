package app_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
	"github.com/pkg/errors"
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
	markup    app.Markup
}

func newElementWithComponent() *elementWithComponent {
	factory := app.NewFactory()
	factory.RegisterComponent(&app.ValidCompo{})

	return &elementWithComponent{
		id:        uuid.New(),
		factory:   factory,
		markup:    html.NewMarkup(factory),
		lastFocus: time.Now(),
	}
}

func (e *elementWithComponent) ID() uuid.UUID {
	return e.id
}

func (e *elementWithComponent) Load(rawurl string) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	compo, err := e.factory.NewComponent(app.ComponentNameFromURL(u))
	if err != nil {
		return err
	}

	if _, err = e.markup.Mount(compo); err != nil {
		return errors.Wrapf(err, "loading %s in test elem %p failed", u, e)
	}
	return nil
}

func (e *elementWithComponent) Contains(c app.Component) bool {
	return e.markup.Contains(c)
}

func (e *elementWithComponent) Render(c app.Component) error {
	_, err := e.markup.Update(c)
	return err
}

func (e *elementWithComponent) LastFocus() time.Time {
	return e.lastFocus
}

func TestElementDB(t *testing.T) {
	tests := []struct {
		scenario string
		capacity int
		function func(t *testing.T, db app.ElementDB)
	}{
		{
			scenario: "should add an element",
			capacity: 1,
			function: testElementDBAdd,
		},
		{
			scenario: "should add an element with components",
			capacity: 1,
			function: testElementDBAddElementWithComponent,
		},
		// {
		// 	scenario: "should fail to add an element when full",
		// 	function: testElementDBAddWhenFull,
		// },
		// {
		// 	scenario: "add element with same id should fail",
		// 	function: testElementDBAddElementWithSameID,
		// },
		// {
		// 	scenario: "should remove an element",
		// 	function: testElementDBRemove,
		// },
		// {
		// 	scenario: "should get an element",
		// 	function: testElementDBElement,
		// },
		// {
		// 	scenario: "should not get an element",
		// 	function: testElementDBElementNotFound,
		// },
		// {
		// 	scenario: "should get an element by component",
		// 	function: testElementDBElementByComponent,
		// },
		// {
		// 	scenario: "should not get an element by component",
		// 	function: testElementDBElementByComponentNotFound,
		// },
		// {
		// 	scenario: "should sort the elements with components",
		// 	function: testElementDBSort,
		// },
		// {
		// 	scenario: "should return the number of elements",
		// 	function: testElementDBLen,
		// },
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t, app.NewElementDB(test.capacity))
		})
	}
}

func testElementDBAdd(t *testing.T, db app.ElementDB) {
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

func testElementDBAddElementWithComponent(t *testing.T, db app.ElementDB) {
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

// func testElementDBAddElementWithSameID(t *testing.T) {
// 	elemDB := newElementDB(42)
// 	elem := newElementWithComponent()

// 	if err := elemDB.Add(elem); err != nil {
// 		t.Fatal(err)
// 	}

// 	err := elemDB.Add(elem)
// 	if err == nil {
// 		t.Fatal("should not add a same element twice")
// 	}
// 	t.Log()

// }

// func testElementDBAddWhenFull(t *testing.T) {
// 	elemDB := newElementDB(42)

// 	newElem := func() Element {
// 		return newElement()
// 	}

// 	for i := 0; i < elemDB.capacity; i++ {
// 		if err := elemDB.Add(newElem()); err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	err := elemDB.Add(newElem())
// 	if err == nil {
// 		t.Fatal("adding an element should return an error")
// 	}
// 	t.Log(err)
// }

// func testElementDBRemove(t *testing.T) {
// 	elemDB := newElementDB(42)
// 	elem := newElementWithComponent()

// 	if err := elemDB.Add(elem); err != nil {
// 		t.Fatal(err)
// 	}

// 	elemDB.Remove(elem)

// 	if l := len(elemDB.elements); l != 0 {
// 		t.Error("elemDB should not have elements:", l)
// 	}
// 	if l := len(elemDB.elementsWithComponents); l != 0 {
// 		t.Error("elemDB should not have elements with components:", l)
// 	}
// }

// func testElementDBElement(t *testing.T) {
// 	elemDB := newElementDB(42)
// 	elem := newElementWithComponent()

// 	if err := elemDB.Add(elem); err != nil {
// 		t.Fatal(err)
// 	}

// 	elemret, ok := elemDB.Element(elem.ID())
// 	if !ok {
// 		t.Fatalf("no element with id %v found", elem.ID())
// 	}
// 	if elemret != elem {
// 		t.Fatal("returned element should be the added element")
// 	}
// }

// func testElementDBElementNotFound(t *testing.T) {
// 	elemDB := newElementDB(42)
// 	if _, ok := elemDB.Element(uuid.New()); ok {
// 		t.Fatal("no element should have been found")
// 	}
// }

// func testElementDBElementByComponent(t *testing.T) {
// 	elem := newElementWithComponent()

// 	compo := &BasicComponent{}
// 	if _, err := elem.env.Mount(compo); err != nil {
// 		t.Fatal(err)
// 	}

// 	elemDB := newElementDB(42)
// 	if err := elemDB.Add(elem); err != nil {
// 		t.Fatal(err)
// 	}

// 	elemret, err := elemDB.ElementByComponent(compo)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if elemret != elem {
// 		t.Fatal("returned element should be the added element")
// 	}
// }

// func testElementDBElementByComponentNotFound(t *testing.T) {
// 	elemDB := newElementDB(42)

// 	if _, err := elemDB.ElementByComponent(&BasicComponent{}); err == nil {
// 		t.Fatal("no element should have been found")
// 	}
// }

// func testElementDBSort(t *testing.T) {
// 	elemDB := newElementDB(42)

// 	for i := 0; i < 10; i++ {
// 		if err := elemDB.Add(newElementWithComponent()); err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	lastElem := newElementWithComponent()
// 	if err := elemDB.Add(lastElem); err != nil {
// 		t.Fatal(err)
// 	}

// 	elems := elemDB.elementsWithComponents
// 	for i, elem := range elems {
// 		if elem.ID() == lastElem.ID() {
// 			elems[i], elems[5] = elems[5], elems[i]
// 			break
// 		}
// 	}

// 	elemDB.Sort()

// 	if elem := elemDB.elementsWithComponents[0]; elem != lastElem {
// 		t.Fatalf("1st element with components should be the last added element: %T", elem)
// 	}
// }

// func testElementDBLen(t *testing.T) {
// 	elemDB := newElementDB(42)

// 	for i := 0; i < 10; i++ {
// 		if err := elemDB.Add(newElementWithComponent()); err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	if l := elemDB.Len(); l != 10 {
// 		t.Fatal("elemDB should have 10 elements:", l)
// 	}
// }
