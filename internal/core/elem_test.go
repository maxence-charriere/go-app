package core

import (
	"testing"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/stretchr/testify/assert"
)

func TestElem(t *testing.T) {
	e := &Elem{}
	e.WhenWindow(func(app.Window) {})
	e.WhenPage(func(app.Page) {})
	e.WhenNavigator(func(app.Navigator) {})
	e.WhenMenu(func(app.Menu) {})
	e.WhenDockTile(func(app.DockTile) {})
	e.WhenStatusMenu(func(app.StatusMenu) {})
	e.WhenNotSet(func() {
		t.Error("WhenNotSet called")
	})

	assert.NoError(t, e.Render(&compo{}))

	e.notSet = true
	e.WhenNotSet(func() {
		t.Log("not set")
	})

	assert.False(t, e.Contains(&compo{}))
	assert.Equal(t, string{}, e.ID())
	assert.Error(t, e.Render(&compo{}))
}

type elem struct {
	Elem
	id string
}

func (e *elem) ID() string {
	return e.id
}

type elemWithCompo struct {
	Elem
	id    string
	compo app.Compo
}

func (e *elemWithCompo) ID() string {
	return e.id
}

func (e *elemWithCompo) Contains(c app.Compo) bool {
	return c != nil && c == e.compo
}

func (e *elemWithCompo) Render(app.Compo) error {
	return nil
}

type compo app.ZeroCompo

func (c *compo) Render() string {
	return `<p></p>`
}

func TestElemDB(t *testing.T) {
	db := NewElemDB()

	// Simple element.
	e := &elem{
		id: uuid.New(),
	}

	db.Put(e)

	e2 := db.GetByID(e.ID())
	assert.False(t, e2.IsNotSet())
	assert.Equal(t, e, e2)

	db.Delete(e)
	e3 := db.GetByID(e.ID())
	assert.True(t, e3.IsNotSet())

	// Element with components.
	ec := &elemWithCompo{
		id: uuid.New(),
	}

	db.Put(ec)
	db.Put(ec)

	c := &compo{}
	ec2 := db.GetByCompo(c)
	assert.True(t, ec2.IsNotSet())

	ec.compo = c
	ec3 := db.GetByCompo(c)
	assert.False(t, ec3.IsNotSet())
	assert.Equal(t, ec, ec3)

	db.Delete(ec)
	ec4 := db.GetByCompo(c)
	assert.True(t, ec4.IsNotSet())
}
