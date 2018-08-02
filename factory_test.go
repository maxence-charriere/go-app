package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	tests := []struct {
		scenario string
		compo    app.Compo
		name     string
		err      bool
	}{
		{
			scenario: "registering and creating a component",
			compo:    &tests.Bar{},
			name:     "tests.bar",
		},
		{
			scenario: "registering a non pointer component",
			compo:    tests.NoPointerCompo{},
			name:     "tests.nopointercompo",
			err:      true,
		},
		{
			scenario: "registering a non pointer to struct component",
			compo: func() *tests.IntCompo {
				intc := tests.IntCompo(42)
				return &intc
			}(),
			name: "tests.intcompo",
			err:  true,
		},
		{
			scenario: "registering a component with no field",
			compo:    &tests.EmptyCompo{},
			name:     "tests.emptycompo",
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			f := app.NewFactory()

			name, err := f.RegisterCompo(test.compo)

			if test.err {
				assert.Error(t, err)
				assert.False(t, f.IsCompoRegistered(test.name))

				_, err = f.NewCompo(test.name)
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.name, name)
			assert.True(t, f.IsCompoRegistered(test.name))

			var c app.Compo
			c, err = f.NewCompo(test.name)

			require.NoError(t, err)
			assert.IsType(t, test.compo, c)
		})
	}
}
