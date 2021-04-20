package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngineInit(t *testing.T) {
	e := engine{}
	e.init()

	assert.NotZero(t, e.UpdateRate)
	assert.NotNil(t, e.Page)
	assert.NotNil(t, e.LocalStorage)
	assert.NotNil(t, e.SessionStorage)
	assert.NotNil(t, e.ResolveStaticResources)
	assert.NotNil(t, e.Body)
	assert.NotNil(t, e.events)
	assert.NotNil(t, e.updates)
	assert.NotNil(t, e.updateQueue)
	assert.NotNil(t, e.defers)
}

func TestEngineScheduleComponentUpdate(t *testing.T) {
	e := engine{}
	e.init()

	h := &hello{}
	e.Mount(h)
	e.Consume()
	require.Empty(t, e.events)
	require.Empty(t, e.updates)
	require.Empty(t, e.updateQueue)

	e.scheduleComponentUpdate(h)
	require.Len(t, e.updates, 1)
	require.Len(t, e.updateQueue, 1)
	require.Equal(t, struct{}{}, e.updates[h])
	require.Equal(t, updateDescriptor{
		compo:    h,
		priority: 2,
	}, e.updateQueue[0])

	e.scheduleComponentUpdate(h)
	require.Len(t, e.updates, 1)
	require.Len(t, e.updateQueue, 1)
}

func TestEngineScheduleNestedComponentUpdate(t *testing.T) {
	e := engine{}
	e.init()

	h := &hello{}
	div := Div().Body(h)

	e.Mount(div)
	e.Consume()
	require.Empty(t, e.events)
	require.Empty(t, e.updates)
	require.Empty(t, e.updateQueue)

	e.scheduleComponentUpdate(h)
	require.Len(t, e.updates, 1)
	require.Len(t, e.updateQueue, 1)
	require.Equal(t, struct{}{}, e.updates[h])
	require.Equal(t, updateDescriptor{
		compo:    h,
		priority: 3,
	}, e.updateQueue[0])
}

func TestEngineUpdateCoponents(t *testing.T) {
	e := engine{}
	e.init()

	foo := &foo{Bar: "bar"}
	e.Mount(foo)
	e.Consume()
	require.Empty(t, e.events)
	require.Empty(t, e.updates)
	require.Empty(t, e.updateQueue)
	bar := foo.root.(*bar)

	e.scheduleComponentUpdate(foo)
	e.scheduleComponentUpdate(bar)
	require.Len(t, e.updates, 2)
	require.Len(t, e.updateQueue, 2)

	e.updateComponents()
	require.Empty(t, e.updates)
	require.Empty(t, e.updateQueue)

	e.updateComponents()
	require.Empty(t, e.updates)
	require.Empty(t, e.updateQueue)
}

func TestEngineExecDeferableEvents(t *testing.T) {
	e := engine{}
	e.init()

	h := &hello{}
	e.Mount(h)
	e.Consume()
	require.Empty(t, e.events)
	require.Empty(t, e.updates)
	require.Empty(t, e.updateQueue)
	require.Empty(t, e.defers)

	called := false

	e.defers = append(e.defers, event{
		source:    h,
		deferable: true,
		function: func(Context) {
			called = true
		},
	})
	require.Len(t, e.defers, 1)

	e.execDeferableEvents()
	require.True(t, called)
	require.Empty(t, e.defers)
}
