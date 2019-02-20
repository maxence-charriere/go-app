package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	tests := []struct {
		scenario string
		compo    Compo
		name     string
		err      bool
	}{
		{
			scenario: "registering and creating a component",
			compo:    &Bar{},
			name:     "app.bar",
		},
		{
			scenario: "registering a non pointer component",
			compo:    NoPointerCompo{},
			name:     "app.nopointercompo",
			err:      true,
		},
		{
			scenario: "registering a non pointer to struct component",
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
			f := newCompoBuilder()

			name, err := f.register(test.compo)

			if test.err {
				assert.Error(t, err)
				assert.False(t, f.isRegistered(test.name))

				_, err = f.new(test.name)
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.name, name)
			assert.True(t, f.isRegistered(test.name))

			var c Compo
			c, err = f.new(test.name)

			require.NoError(t, err)
			assert.IsType(t, test.compo, c)
		})
	}
}
