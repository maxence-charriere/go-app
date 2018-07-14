package tests

import (
	"reflect"
	"testing"

	"github.com/murlokswarm/app"
)

// TestFactory is a test suite used to ensure that all the factory
// implementations behave the same.
func TestFactory(t *testing.T, setup func() app.Factory) {
	tests := []struct {
		scenario string
		compo    app.Compo
		name     string
		err      bool
	}{
		{
			scenario: "registering and creating a component",
			compo:    &Bar{},
			name:     "tests.bar",
		},
		{
			scenario: "registering a non pointer component",
			compo:    NoPointerCompo{},
			name:     "tests.nopointercompo",
			err:      true,
		},
		{
			scenario: "registering a non pointer to struct component",
			compo: func() *IntCompo {
				intc := IntCompo(42)
				return &intc
			}(),
			name: "tests.intcompo",
			err:  true,
		},
		{
			scenario: "registering a component with no field",
			compo:    &EmptyCompo{},
			name:     "tests.emptycompo",
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			factory := setup()

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

			var newCompo app.Compo

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
