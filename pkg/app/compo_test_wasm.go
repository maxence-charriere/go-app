package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type boo struct {
	Compo
	Dismount func()
}

func (b *boo) OnDismount() {
	if b.Dismount != nil {
		b.Dismount()
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
		Dismount: func() {
			called = true
		},
	}

	mount(c)
	c.dismount()
	require.True(t, called)
}
