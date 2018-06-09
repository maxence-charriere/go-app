package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type control struct {
	id    string
	compo bool
}

func (e *control) ID() string {
	return e.id
}

func (e *control) Contains(c Component) bool {
	return e.compo
}

func TestControlDB(t *testing.T) {
	tests := []struct {
		scenario string
		elem     *control
		byID     string
		errID    bool
		errCompo bool
	}{
		{
			scenario: "elem with component",
			elem: &control{
				id:    "hello",
				compo: true,
			},
			byID: "hello",
		},
		{
			scenario: "elem without component",
			elem: &control{
				id: "hello",
			},
			byID:     "hello",
			errCompo: true,
		},
		{
			scenario: "no elem",
			errID:    true,
			errCompo: true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			controlDB := NewControlDB()

			if test.elem != nil {
				controlDB.Add(test.elem)

				defer func() {
					controlDB.Remove(test.elem)
					_, err := controlDB.ControlByID(test.elem.ID())
					assert.Error(t, err)
				}()
			}

			elem, err := controlDB.ControlByID(test.byID)
			if test.errID {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.elem.ID(), elem.ID())
			}

			elem, err = controlDB.ControlByComponent(&Bar{})
			if test.errCompo {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.elem.ID(), elem.ID())
			}
		})
	}
}
