package html

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
)

func TestEncoder(t *testing.T) {
	factory := app.NewFactory()
	factory.Register(&tests.Hello{})
	factory.Register(&tests.World{})

	tests := []struct {
		scenario string
		function func(t *testing.T, markup *Markup)
	}{
		{
			scenario: "encoding a component",
			function: testEncoderEncode,
		},
		{
			scenario: "encoding a zero tag returns an error",
			function: testEncoderEncodeZeroTag,
		},
		{
			scenario: "encoding a tag with a zero tag child returns an error",
			function: testEncoderEncodeChildZeroTag,
		},
		{
			scenario: "encoding a not mounted component tag returns an error",
			function: testEncoderEncodeNotMountedCompo,
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
	root, err := markup.Mount(&tests.Hello{
		Name: "Maxence",
	})
	if err != nil {
		t.Fatal(err)
	}

	buff := &bytes.Buffer{}
	enc := NewEncoder(buff, markup, true)

	if err = enc.Encode(root); err != nil {
		t.Fatal(err)
	}
	t.Log(buff.String())
}

func testEncoderEncodeZeroTag(t *testing.T, markup *Markup) {
	root := app.Tag{}

	buff := &bytes.Buffer{}
	enc := NewEncoder(buff, markup, true)

	err := enc.Encode(root)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testEncoderEncodeChildZeroTag(t *testing.T, markup *Markup) {
	root, err := markup.Mount(&tests.Hello{
		Name: "Maxence",
	})
	if err != nil {
		t.Fatal(err)
	}

	root.Children[0].Type = app.ZeroTag

	buff := &bytes.Buffer{}
	enc := NewEncoder(buff, markup, true)

	if err = enc.Encode(root); err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testEncoderEncodeNotMountedCompo(t *testing.T, markup *Markup) {
	root := app.Tag{
		Name: "html.world",
		ID:   uuid.New(),
		Type: app.CompoTag,
	}

	buff := &bytes.Buffer{}
	enc := NewEncoder(buff, markup, true)

	err := enc.Encode(root)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func BenchmarkEncoder(b *testing.B) {
	factory := app.NewFactory()
	factory.Register(&tests.Hello{})
	factory.Register(&tests.World{})

	markup := NewMarkup(factory)

	hello := &tests.Hello{
		Name: "JonhyMaxoo",
	}

	root, _ := markup.Mount(hello)

	for i := 0; i < b.N; i++ {
		var v bytes.Buffer
		enc := NewEncoder(&v, markup, true)
		enc.Encode(root)
	}
}

func TestAttrValueFormatter(t *testing.T) {
	factory := app.NewFactory()
	factory.Register(&tests.World{})
	compoID := uuid.New()

	tests := []struct {
		scenario  string
		formatter AttrValueFormatter
		expected  string
	}{
		{
			scenario: "href no format",
			formatter: AttrValueFormatter{
				Name:       "href",
				Value:      "tests.world",
				FormatHref: false,
				Factory:    factory,
			},
			expected: "tests.world",
		},
		{
			scenario: "href with format",
			formatter: AttrValueFormatter{
				Name:       "href",
				Value:      "tests.world",
				FormatHref: true,
				Factory:    factory,
			},
			expected: "compo:///tests.world",
		},
		{
			scenario: "html handler",
			formatter: AttrValueFormatter{
				Name:    "onclick",
				Value:   "OnTest",
				CompoID: compoID,
			},
			expected: fmt.Sprintf(
				`callGoEventHandler('%s', 'OnTest', this, event)`,
				compoID,
			),
		},
		{
			scenario: "html js handler",
			formatter: AttrValueFormatter{
				Name:    "onclick",
				Value:   "js:alert('test')",
				CompoID: compoID,
			},
			expected: "alert('test')",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			if value := test.formatter.Format(); value != test.expected {
				t.Error("expected:", test.expected)
				t.Error("value   :", value)
			}
		})
	}
}
