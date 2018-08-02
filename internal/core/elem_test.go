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
	e.WhenErr(func(err error) {
		t.Error("WhenErr called:", err)
	})

	e.SetErr(app.ErrElemNotSet)
	e.WhenErr(func(err error) {
		t.Log("WhenErr called:", err)
	})

	assert.Equal(t, "", e.ID())
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

type compo app.ZeroCompo

func (c *compo) Render() string {
	return `<p></p>`
}

func TestElemDB(t *testing.T) {
	db := NewElemDB()

	// Simple element.
	e := &elem{
		id: uuid.New().String(),
	}

	db.Put(e)

	e2 := db.GetByID(e.ID())
	assert.NoError(t, e2.Err())
	assert.Equal(t, e, e2)

	db.Delete(e)
	e3 := db.GetByID(e.ID())
	assert.Error(t, e3.Err())

	// Element with components.
	ec := &elemWithCompo{
		id: uuid.New().String(),
	}

	db.Put(ec)
	db.Put(ec)

	c := &compo{}
	ec2 := db.GetByCompo(c)
	assert.Error(t, ec2.Err())

	ec.compo = c
	ec3 := db.GetByCompo(c)
	assert.NoError(t, ec3.Err())
	assert.Equal(t, ec, ec3)

	db.Delete(ec)
	ec4 := db.GetByCompo(c)
	assert.Error(t, ec4.Err())
}
