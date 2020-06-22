package app

import "testing"

func TestRange(t *testing.T) {
	testUpdate(t, []updateTest{
		{
			scenario: "range slice is updated",
			a: Div().Body(
				Range([]string{"hello", "world"}).Slice(func(i int) UI {
					src := []string{"hello", "world"}
					return Text(src[i])
				}),
			),
			b: Div().Body(
				Range([]string{"hello", "maxoo"}).Slice(func(i int) UI {
					src := []string{"hello", "maxoo"}
					return Text(src[i])
				}),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Text("hello"),
				},
				{
					Path:     TestPath(1),
					Expected: Text("maxoo"),
				},
			},
		},
		{
			scenario: "range slice is updated to be empty",
			a: Div().Body(
				Range([]string{"hello", "world"}).Slice(func(i int) UI {
					src := []string{"hello", "world"}
					return Text(src[i])
				}),
			),
			b: Div().Body(
				Range([]string{}).Slice(func(i int) UI {
					src := []string{"hello", "maxoo"}
					return Text(src[i])
				}),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: nil,
				},
				{
					Path:     TestPath(1),
					Expected: nil,
				},
			},
		},
		{
			scenario: "range map is updated",
			a: Div().Body(
				Range(map[string]string{"key": "value"}).Map(func(k string) UI {
					src := map[string]string{"key": "value"}
					return Text(src[k])
				}),
			),
			b: Div().Body(
				Range(map[string]string{"key": "value"}).Map(func(k string) UI {
					src := map[string]string{"key": "maxoo"}
					return Text(src[k])
				}),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Text("maxoo"),
				},
			},
		},
	})
}
