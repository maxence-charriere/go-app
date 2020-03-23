package app

import "testing"

type boo struct {
	Compo
}

func (b *boo) Render() UI {
	return Text("foo")
}

type booWithDefaultRender struct {
	Compo
}

func TestCompoUmountedUpdate(t *testing.T) {
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
