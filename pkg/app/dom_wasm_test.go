package app

import (
	"syscall/js"
	"testing"

	"github.com/stretchr/testify/require"
)

type ctxmenu ZeroCompo

func (c *ctxmenu) Render() string {
	return `<div></div>`
}

func TestDom(t *testing.T) {
	tests := []struct {
		scenario     string
		imports      []Compo
		ctxmenu      Compo
		compo        Compo
		modifier     func(Compo)
		expectedRoot *node
		compoCount   int
	}{
		// {
		// 	scenario: "simple component",
		// 	imports: []Compo{
		// 		&Bar{},
		// 		&ctxmenu{},
		// 	},
		// 	ctxmenu: &ctxmenu{},
		// 	compo:   &Bar{},
		// 	expectedRoot: &node{
		// 		name:      "h2",
		// 		compoName: "app.bar",
		// 		children: []*node{
		// 			{
		// 				text: "Bar",
		// 			},
		// 		},
		// 	},
		// 	compoCount: 2,
		// },
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			b := make(compoBuilder)
			b.imports(test.ctxmenu)
			for _, c := range test.imports {
				b.imports(c)
			}

			dom := &dom{
				compoBuilder: b,
				callOnUI: func(f func()) {
					f()
				},
				trackCursorPosition: func(js.Value) {},
				contextMenu:         test.ctxmenu,
			}

			err := dom.newBody(test.compo)
			require.NoError(t, err, "new body")

			if test.modifier != nil {
				test.modifier(test.compo)
				err = dom.render(test.compo)
				require.NoError(t, err, "render")
			}

			compareNode(t, test.expectedRoot, dom.root)
			require.Equal(t, test.compoCount, len(dom.components), "component count")

			dom.clean()
			require.Empty(t, dom.components)
			require.Nil(t, dom.root)
			require.Nil(t, dom.ctxMenuRoot)
		})
	}
}

func compareNode(t *testing.T, x, y *node) {
	require.Equal(t, x.Name, y.Name)
	require.Equal(t, x.Text, y.Text)
	require.Equal(t, x.Attrs, y.Attrs)
	require.Equal(t, x.CompoName, y.CompoName)
	require.Equal(t, x.isEnd, y.isEnd)

	require.Equal(t, len(x.Children), len(x.Children), "children len")
	for i := 0; i < len(x.Children); i++ {
		compareNode(t, x.Children[i], y.Children[i])
	}
}
