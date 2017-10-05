package app

// type BasicComponent ZeroCompo

// func (c *BasicComponent) Render() string {
// 	return "<div><div>"
// }

// type element struct {
// 	id uuid.UUID
// }

// func newElement() *element {
// 	elem := &element{
// 		id: uuid.New(),
// 	}
// 	return elem
// }

// func (e *element) ID() uuid.UUID {
// 	return e.id
// }

// type elementWithComponent struct {
// 	id           uuid.UUID
// 	compoBuilder CompoBuilder
// 	lastFocus    time.Time
// 	env          Env
// }

// func newElementWithComponent() *elementWithComponent {
// 	compoBuilder := NewCompoBuilder()
// 	compoBuilder.Register(&BasicComponent{})

// 	return &elementWithComponent{
// 		id:           uuid.New(),
// 		compoBuilder: compoBuilder,
// 		env:          NewEnv(compoBuilder),
// 		lastFocus:    time.Now(),
// 	}
// }

// func (e *elementWithComponent) ID() uuid.UUID {
// 	return e.id
// }

// func (e *elementWithComponent) Load(rawurl string) error {
// 	u, err := url.Parse(rawurl)
// 	if err != nil {
// 		return err
// 	}

// 	componame, ok := ComponentNameFromURL(u)
// 	if !ok {
// 		return nil
// 	}

// 	compo, err := e.compoBuilder.New(componame)
// 	if err != nil {
// 		return err
// 	}

// 	if _, err = e.env.Mount(compo); err != nil {
// 		return errors.Wrapf(err, "loading %s in test elem %p failed", u, e)
// 	}
// 	return nil
// }

// func (e *elementWithComponent) Contains(c Component) bool {
// 	return e.env.Contains(c)
// }

// func (e *elementWithComponent) Render(c Component) error {
// 	_, err := e.env.Update(c)
// 	return err
// }

// func (e *elementWithComponent) LastFocus() time.Time {
// 	return e.lastFocus
// }

// func TestElementDB(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		test func(t *testing.T)
// 	}{
// 		{
// 			name: "should add an element",
// 			test: testElementDBAdd,
// 		},
// 		{
// 			name: "should add an element with components",
// 			test: testElementDBAddElementWithComponent,
// 		},
// 		{
// 			name: "should fail to add an element when full",
// 			test: testElementDBAddWhenFull,
// 		},
// 		{
// 			name: "add element with same id should fail",
// 			test: testElementDBAddElementWithSameID,
// 		},
// 		{
// 			name: "should remove an element",
// 			test: testElementDBRemove,
// 		},
// 		{
// 			name: "should get an element",
// 			test: testElementDBElement,
// 		},
// 		{
// 			name: "should not get an element",
// 			test: testElementDBElementNotFound,
// 		},
// 		{
// 			name: "should get an element by component",
// 			test: testElementDBElementByComponent,
// 		},
// 		{
// 			name: "should not get an element by component",
// 			test: testElementDBElementByComponentNotFound,
// 		},
// 		{
// 			name: "should sort the elements with components",
// 			test: testElementDBSort,
// 		},
// 		{
// 			name: "should return the number of elements",
// 			test: testElementDBLen,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, test.test)
// 	}
// }

// func testElementDBAdd(t *testing.T) {
// 	elemDB := newElementDB(42)

// 	if err := elemDB.Add(newElement()); err != nil {
// 		t.Fatal(err)
// 	}

// 	if l := len(elemDB.elements); l != 1 {
// 		t.Error("elemDB should have 1 element:", l)
// 	}
// 	if l := len(elemDB.elementsWithComponents); l != 0 {
// 		t.Error("elemDB should not have an element with components")
// 	}
// }

// func testElementDBAddElementWithComponent(t *testing.T) {
// 	elemDB := newElementDB(42)

// 	if err := elemDB.Add(newElementWithComponent()); err != nil {
// 		t.Fatal(err)
// 	}

// 	if l := len(elemDB.elements); l != 1 {
// 		t.Error("elemDB should have 1 element:", l)
// 	}
// 	if l := len(elemDB.elementsWithComponents); l != 1 {
// 		t.Error("elemDB should have 1 element with components:", l)
// 	}
// }

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
