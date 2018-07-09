package html

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/murlokswarm/app"
)

type CompoWithFields struct {
	app.ZeroCompo
	secret             string
	funcHandler        func()
	funcWithArgHandler func(int)

	String     string
	Bool       bool
	NotSetBool bool
	Int        int
	Uint       uint
	Float      float64
	Struct     struct {
		A int
		B string
	}
	Time time.Time
}

func (c *CompoWithFields) Render() string {
	return `
<div>
	<div>String: {{.String}}</div>
	<div>raw String: {{raw .String}}</div>
	<div>Bool: {{.Bool}}</div>
	<div>Int: {{.Int}}</div>
	<div>Uint: {{.Uint}}</div>
	<div>Float: {{.Float}}</div>
	<div>Struct: {{.Struct}}</div>
	<html.compo obj="{{json .Struct}}">	
	<div>Time: {{time .Time "2006"}}</div>
	<div>{{hello .String}}</div>
	<div>compo String: {{compo "html.compo"}}</div>	
</div>
	`
}

func (c *CompoWithFields) Funcs() map[string]interface{} {
	return map[string]interface{}{
		"hello": func(string) string { return "hello" },
	}
}

type CompoBadTemplate app.ZeroCompo

func (c *CompoBadTemplate) Render() string {
	return `{{}}`
}

type CompoBadTemplate2 app.ZeroCompo

func (c *CompoBadTemplate2) Render() string {
	return `{{.Bye}}`
}

func TestDecodeComponent(t *testing.T) {
	s := struct {
		A int
		B string
	}{
		A: 42,
		B: "foobar",
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	sjson := string(data)

	n, err := decodeComponent(&CompoWithFields{
		String: "<br>",
		Time:   time.Now(),
		Struct: s,
	})
	require.NoError(t, err)

	root := n.(*elemNode)
	raw := root.children[1].(*elemNode).children[1].(*elemNode)
	assert.Equal(t, "br", raw.TagName())

	compo := root.children[7].(*compoNode)
	assert.Equal(t, sjson, compo.fields["obj"])

	year := strconv.Itoa(time.Now().Year())
	timetext := root.children[8].(*elemNode).children[0].(*textNode)
	assert.Equal(t, "Time: "+year, timetext.Text())

	hello := root.children[9].(*elemNode).children[0].(*textNode)
	assert.Equal(t, "hello", hello.Text())

	_, err = decodeComponent(&CompoBadTemplate{})
	assert.Error(t, err)

	_, err = decodeComponent(&CompoBadTemplate2{})
	assert.Error(t, err)
}

func TestMapComponentFields(t *testing.T) {
	tests := []struct {
		scenario string
		attrs    map[string]string
		expected CompoWithFields
		err      bool
	}{
		{
			scenario: "skip mapping nil",
			attrs:    nil,
		},
		{
			scenario: "skip mapping an anonymous field",
			attrs:    map[string]string{"zerocompo": `{"placeholder": 42}`},
		},
		{
			scenario: "skip mapping an unexported field",
			attrs:    map[string]string{"secret": "pandore"},
		},
		{
			scenario: "map a string",
			attrs:    map[string]string{"string": "hello"},
			expected: CompoWithFields{
				String: "hello",
			},
		},
		{
			scenario: "map a bool",
			attrs:    map[string]string{"bool": "true"},
			expected: CompoWithFields{
				Bool: true,
			},
		},
		{
			scenario: "map a naked bool",
			attrs:    map[string]string{"bool": ""},
			expected: CompoWithFields{
				Bool: true,
			},
		},
		{
			scenario: "map a non boolean value to bool returns an error",
			attrs:    map[string]string{"bool": "lolilol"},
			err:      true,
		},
		{
			scenario: "map an int",
			attrs:    map[string]string{"int": "-42"},
			expected: CompoWithFields{
				Int: -42,
			},
		},
		{
			scenario: "map a non int value to int returns an error",
			attrs:    map[string]string{"int": "lolilol"},
			err:      true,
		},
		{
			scenario: "map an uint",
			attrs:    map[string]string{"uint": "21"},
			expected: CompoWithFields{
				Uint: 21,
			},
		},
		{
			scenario: "map a non uint value to uint returns an error",
			attrs:    map[string]string{"uint": "lolilol"},
			err:      true,
		},
		{
			scenario: "map a float",
			attrs:    map[string]string{"float": "42.42"},
			expected: CompoWithFields{
				Float: 42.42,
			},
		},
		{
			scenario: "map a non float value to float returns an error",
			attrs:    map[string]string{"float": "42.world"},
			err:      true,
		},
		{
			scenario: "map a struct",
			attrs:    map[string]string{"struct": `{"A": 42, "B": "world"}`},
			expected: CompoWithFields{
				Struct: struct {
					A int
					B string
				}{
					A: 42,
					B: "world",
				},
			},
		},
		{
			scenario: "map a struct with invalid fields returns an error",
			attrs:    map[string]string{"struct": `{"A": "world", "B": 42}`},
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			var c CompoWithFields

			err := mapComponentFields(&c, test.attrs)
			if test.err {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, test.expected, c)
		})
	}
}
