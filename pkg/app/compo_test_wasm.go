package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type boo struct {
	Compo

	Value int

	onDismount func()
	onUpdate   func()
}

func (b *boo) OnDismount() {
	if b.onDismount != nil {
		b.onDismount()
	}
}

func (b *boo) OnUpdate() {
	if b.onUpdate != nil {
		b.onUpdate()
	}
}

func (b *boo) Render() UI {
	return Text("foo")
}

type booWithDefaultRender struct {
	Compo
}

func TestCompoUnmountedUpdate(t *testing.T) {
	tests := []struct {
		scenario string
		compo    Composer
	}{
		{
			scenario: "component with redefined render is updated",
			compo:    &boo{},
		},
		{
			scenario: "component without redefined render is updated",
			compo:    &booWithDefaultRender{},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			dispatcher = func(f func()) {
				f()
			}
			defer func() {
				dispatcher = Dispatch
			}()

			test.compo.Update()
		})
	}
}

func TestCompoDismount(t *testing.T) {
	called := false

	c := &boo{
		onDismount: func() {
			called = true
		},
	}

	mount(c)
	c.dismount()
	require.True(t, called)
}

func TestCompoUpdatable(t *testing.T) {
	called := false
	onUpdate := func() {
		called = true
	}

	a := &boo{Value: 42, onUpdate: onUpdate}
	err := mount(a)
	require.NoError(t, err)

	b := &boo{Value: 42, onUpdate: onUpdate}
	err = update(a, b)
	require.NoError(t, err)
	require.False(t, called)

	c := &boo{Value: 21, onUpdate: onUpdate}
	err = update(a, c)
	require.NoError(t, err)
	require.Equal(t, 21, a.Value)
	require.True(t, called)
}
