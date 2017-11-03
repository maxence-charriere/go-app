package html

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

func TestEncoder(t *testing.T) {
	factory := app.NewFactory()
	factory.RegisterComponent(&Hello{})
	factory.RegisterComponent(&World{})

	tests := []struct {
		scenario string
		function func(t *testing.T, markup *Markup)
	}{
		{
			scenario: "should encode component",
			function: testEncoderEncode,
		},
		{
			scenario: "encode a zero tag should fail",
			function: testEncoderEncodeZeroTag,
		},
		{
			scenario: "encode tag with a zero tag child should fail",
			function: testEncoderEncodeChildZeroTag,
		},
		{
			scenario: "encode not mounted component tag should fail",
			function: testEncoderEncodeNotMountedComponent,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			markup := NewMarkup(factory)
			test.function(t, markup)
		})
	}
}

func testEncoderEncode(t *testing.T, markup *Markup) {
	root, err := markup.Mount(&Hello{
		Name: "Maxence",
	})
	if err != nil {
		t.Fatal(err)
	}

	buff := &bytes.Buffer{}
	enc := NewEncoder(buff, markup)

	if err = enc.Encode(root); err != nil {
		t.Fatal(err)
	}
	t.Log(buff.String())
}

func testEncoderEncodeZeroTag(t *testing.T, markup *Markup) {
	root := app.Tag{}

	buff := &bytes.Buffer{}
	enc := NewEncoder(buff, markup)

	err := enc.Encode(root)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testEncoderEncodeChildZeroTag(t *testing.T, markup *Markup) {
	root, err := markup.Mount(&Hello{
		Name: "Maxence",
	})
	if err != nil {
		t.Fatal(err)
	}

	root.Children[0].Type = app.ZeroTag

	buff := &bytes.Buffer{}
	enc := NewEncoder(buff, markup)

	if err = enc.Encode(root); err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testEncoderEncodeNotMountedComponent(t *testing.T, markup *Markup) {
	root := app.Tag{
		Name: "html.world",
		ID:   uuid.New(),
		Type: app.CompoTag,
	}

	buff := &bytes.Buffer{}
	enc := NewEncoder(buff, markup)

	err := enc.Encode(root)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func BenchmarkEncoder(b *testing.B) {
	factory := app.NewFactory()
	factory.RegisterComponent(&Hello{})
	factory.RegisterComponent(&World{})

	markup := NewMarkup(factory)

	hello := &Hello{
		Name: "JonhyMaxoo",
	}

	root, _ := markup.Mount(hello)

	for i := 0; i < b.N; i++ {
		var v bytes.Buffer
		enc := NewEncoder(&v, markup)
		enc.Encode(root)
	}
}
