package tests

import (
	"testing"

	"github.com/murlokswarm/app"
	// "github.com/murlokswarm/app/internal/dom"
	"github.com/murlokswarm/app/internal/dom.v2"
)

type BenchCompo struct {
	N int
}

func (b *BenchCompo) Render() string {
	return `
	<div>
		<!-- Comment -->	
		<h1>Iteration {{.N}}</h1>
		<br>
		<tests.benchsubcompo n="{{.N}}">
		<input type="text" required onchange="Test">
		<svg>
			<path d="M 42.42 Z "></path>
			<path d="M 21.21 Z " />
		</svg>
		<a href="html.Foo"></a>
	</div>`
}

type BenchSubCompo struct {
	N int
}

func (b *BenchSubCompo) Render() string {
	return `
	<p>{{.N}}</p>`
}

// func BenchmarkDom(b *testing.B) {
// 	b.ReportAllocs()

// 	f := app.NewFactory()
// 	f.RegisterCompo(&BenchCompo{})
// 	f.RegisterCompo(&BenchSubCompo{})

// 	d := dom.NewDOM(f, dom.JsToGoHandler, dom.HrefCompoFmt)

// 	c := &BenchCompo{}

// 	for n := 0; n < b.N; n++ {
// 		c.N = n

// 		if n == 0 {
// 			if _, err := d.New(c); err != nil {
// 				b.Fatal(err)
// 			}
// 			continue
// 		}

// 		if _, err := d.Update(c); err != nil {
// 			b.Fatal(err)
// 		}
// 	}
// }

func BenchmarkDom(b *testing.B) {
	b.ReportAllocs()

	f := app.NewFactory()
	f.RegisterCompo(&BenchCompo{})
	f.RegisterCompo(&BenchSubCompo{})

	d := &dom.Engine{
		Factory: f,
		AttrTransforms: []dom.Transform{
			dom.JsToGoHandler,
			dom.HrefCompoFmt,
		},
	}

	c := &BenchCompo{}

	for n := 0; n < b.N; n++ {
		c.N = n

		if n == 0 {
			if err := d.New(c); err != nil {
				b.Fatal(err)
			}
			continue
		}

		if err := d.Render(c); err != nil {
			b.Fatal(err)
		}
	}
}
