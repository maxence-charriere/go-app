package app

import (
	"reflect"
	"testing"
)

// TestFactory is a test suite used to ensure that all the factory
// implementations behave the same.
func TestFactory(t *testing.T) {
	tests := []struct {
		scenario string
		compo    Component
		name     string
		err      bool
	}{
		{
			scenario: "registering and creating a component",
			compo:    &Bar{},
			name:     "app.bar",
		},
		{
			scenario: "registering a non struct component",
			compo: func() *IntCompo {
				intc := IntCompo(42)
				return &intc
			}(),
			name: "app.intcompo",
			err:  true,
		},
		{
			scenario: "registering a component with no field",
			compo:    &EmptyCompo{},
			name:     "app.emptycompo",
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			factory := NewFactory()

			name, err := factory.Register(test.compo)
			if test.err && err == nil {
				t.Error("register succeeded")
			} else if !test.err {
				if err != nil {
					t.Error(err)
				}
				if name != test.name {
					t.Errorf("registered name is not %s: %s", test.name, name)
				}
			}

			if ok := factory.Registered(test.name); test.err && ok {
				t.Error("component is registered")
			} else if !test.err && !ok {
				t.Error("component is not registered")
			}

			var newCompo Component

			newCompo, err = factory.New(test.name)
			if test.err && err == nil {
				t.Error("component is created")
			} else if !test.err {
				if err != nil {
					t.Error(err)
				}

				ctype := reflect.TypeOf(test.compo)
				ntype := reflect.TypeOf(newCompo)

				if ntype != ctype {
					t.Errorf("created component is not %v: %v", ctype, ntype)
				}
			}
		})
	}
}
