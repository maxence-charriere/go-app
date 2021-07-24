package app

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngineInit(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	assert.NotZero(t, e.UpdateRate)
	assert.NotNil(t, e.Page)
	assert.NotNil(t, e.LocalStorage)
	assert.NotNil(t, e.SessionStorage)
	assert.NotNil(t, e.ResolveStaticResources)
	assert.NotNil(t, e.Body)
	assert.NotNil(t, e.dispatches)
	assert.NotNil(t, e.updates)
	assert.NotNil(t, e.updateQueue)
	assert.NotNil(t, e.defers)
}

func TestEngineDispatch(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	e.Dispatch(Dispatch{})

	require.Len(t, e.dispatches, 1)

	d := <-e.dispatches
	require.Equal(t, Update, d.Mode)
	require.Equal(t, e.Body, d.Source)
	require.NotNil(t, d.Function)
}

func TestEngineEmit(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	foo := &foo{Bar: "bar"}
	e.Mount(foo)
	e.Consume()
	require.Empty(t, e.dispatches)
	require.Empty(t, e.updates)
	require.Empty(t, e.updateQueue)

	bar := foo.children()[0].(*bar)

	emitted := false
	e.Emit(bar, func() {
		emitted = true
	})
	require.True(t, emitted)
	require.Len(t, e.dispatches, 1)

	e.Emit(bar, nil)
}

func TestEngineHandleDispatch(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		e := engine{}
		e.init()
		defer e.Close()

		bar := &bar{}
		e.Mount(bar)
		e.Consume()

		called := false
		e.handleDispatch(Dispatch{
			Mode:     Update,
			Source:   bar,
			Function: func(Context) { called = true },
		})
		require.True(t, called)
		require.NotEmpty(t, e.updateQueue)
	})

	t.Run("defer", func(t *testing.T) {
		e := engine{}
		e.init()
		defer e.Close()

		bar := &bar{}
		e.Mount(bar)
		e.Consume()

		called := false
		e.handleDispatch(Dispatch{
			Mode:     Defer,
			Source:   bar,
			Function: func(Context) { called = true },
		})
		require.Empty(t, e.updateQueue)
		require.Len(t, e.defers, 1)
		require.False(t, called)
	})

	t.Run("next", func(t *testing.T) {
		e := engine{}
		e.init()
		defer e.Close()

		bar := &bar{}
		e.Mount(bar)
		e.Consume()

		called := false
		e.handleDispatch(Dispatch{
			Mode:     Next,
			Source:   bar,
			Function: func(Context) { called = true },
		})
		require.True(t, called)
		require.Empty(t, e.updateQueue)
	})
}

func TestEngineScheduleComponentUpdate(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	h := &hello{}
	e.scheduleComponentUpdate(h)
	require.Empty(t, e.updates)
	require.Empty(t, e.updateQueue)

	e.Mount(h)
	e.Consume()
	require.Empty(t, e.dispatches)
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
	defer e.Close()

	h := &hello{}
	div := Div().Body(h)
	e.scheduleComponentUpdate(h)
	require.Empty(t, e.updates)
	require.Empty(t, e.updateQueue)

	e.Mount(div)
	e.Consume()
	require.Empty(t, e.dispatches)
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
	defer e.Close()

	foo := &foo{Bar: "bar"}
	e.Mount(foo)
	e.Consume()
	require.Empty(t, e.dispatches)
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
	defer e.Close()

	h := &hello{}
	e.Mount(h)
	e.Consume()
	require.Empty(t, e.dispatches)
	require.Empty(t, e.updates)
	require.Empty(t, e.updateQueue)
	require.Empty(t, e.defers)

	called := false

	e.defers = append(e.defers, Dispatch{
		Mode:   Defer,
		Source: h,
		Function: func(Context) {
			called = true
		},
	})
	require.Len(t, e.defers, 1)

	e.execDeferableEvents()
	require.True(t, called)
	require.Empty(t, e.defers)
}

func TestEngineHandlePost(t *testing.T) {
	isAppHandleCalled := false
	isHandleACalled := false
	isHandleBCalled := false
	isHandleCCalled := false

	e := engine{
		ActionHandlers: map[string]ActionHandler{
			"/test": func(ctx Context, a Action) {
				isAppHandleCalled = true
			},
		},
	}
	e.init()
	defer e.Close()

	h := &hello{}
	e.Mount(h)
	e.Consume()

	e.Handle("/test", h, func(ctx Context, a Action) {
		isHandleACalled = true
	})

	e.Handle("/test", h, func(ctx Context, a Action) {
		isHandleBCalled = true
	})

	f := &foo{}
	e.Handle("/test", f, func(ctx Context, a Action) {
		isHandleCCalled = true
	})

	e.Post(Action{Name: "/test"})
	e.Consume()

	require.True(t, isAppHandleCalled)
	require.True(t, isHandleACalled)
	require.True(t, isHandleBCalled)
	require.False(t, isHandleCCalled)
}

func TestSortUpdateDescriptors(t *testing.T) {
	utests := []struct {
		scenario string
		in       []updateDescriptor
		out      []updateDescriptor
	}{
		{
			scenario: "nil",
		},
		{
			scenario: "empty",
			in:       []updateDescriptor{},
			out:      []updateDescriptor{},
		},
		{
			scenario: "single value",
			in: []updateDescriptor{
				{priority: 42},
			},
			out: []updateDescriptor{
				{priority: 42},
			},
		},
		{
			scenario: "two values",
			in: []updateDescriptor{
				{priority: 42},
				{priority: 21},
			},
			out: []updateDescriptor{
				{priority: 21},
				{priority: 42},
			},
		},
		{
			scenario: "multiple values",
			in: []updateDescriptor{
				{priority: 43},
				{priority: 2},
				{priority: 9},
				{priority: 36},
				{priority: 21},
				{priority: 198},
				{priority: 9},
				{priority: 1},
			},
			out: []updateDescriptor{
				{priority: 1},
				{priority: 2},
				{priority: 9},
				{priority: 9},
				{priority: 21},
				{priority: 36},
				{priority: 43},
				{priority: 198},
			},
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			sortUpdateDescriptors(u.in)
			require.Equal(t, u.out, u.in)
		})
	}
}

const (
	sortUpdateBenchCount = 1000
)

func BenchmarkSortUpdateDescriptor(b *testing.B) {
	benchSortUpdateDescriptor(b, sortUpdateBenchCount, sortUpdateDescriptors)

}

func BenchmarkSortUpdateDescriptorStd(b *testing.B) {
	benchSortUpdateDescriptor(b, sortUpdateBenchCount, func(d []updateDescriptor) {
		sort.Slice(d, func(a, b int) bool {
			return d[a].priority < d[b].priority
		})
	})
}

func BenchmarkSortUpdateDescriptorOrdered(b *testing.B) {
	benchSortUpdateDescriptorOrdered(b, sortUpdateBenchCount, sortUpdateDescriptors)

}

func BenchmarkSortUpdateDescriptorOrderedStd(b *testing.B) {
	benchSortUpdateDescriptorOrdered(b, sortUpdateBenchCount, func(d []updateDescriptor) {
		sort.Slice(d, func(a, b int) bool {
			return d[a].priority < d[b].priority
		})
	})
}

func benchSortUpdateDescriptor(b *testing.B, n int, sort func([]updateDescriptor)) {
	rand.Seed(time.Now().UnixNano())

	d := make([]updateDescriptor, n)
	for i := range d {
		d[i].compo = &hello{}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		unorderUpdateDescriptor(b, d)
		sort(d)
	}
}

func benchSortUpdateDescriptorOrdered(b *testing.B, n int, sort func([]updateDescriptor)) {
	rand.Seed(time.Now().UnixNano())

	d := make([]updateDescriptor, n)
	for i := range d {
		d[i].compo = &hello{}
	}
	unorderUpdateDescriptor(b, d)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sort(d)
	}
}

func unorderUpdateDescriptor(b *testing.B, d []updateDescriptor) {
	b.StopTimer()
	for i := range d {
		d[i].priority = rand.Intn(1000)
	}
	b.StartTimer()
}
