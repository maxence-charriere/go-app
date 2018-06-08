package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type element struct {
	id    string
	compo bool
}

func (e *element) ID() string {
	return e.id
}

func (e *element) Contains(c Component) bool {
	return e.compo
}

func TestElementDB(t *testing.T) {
	tests := []struct {
		scenario string
		elem     *element
		byID     string
		errID    bool
		errCompo bool
	}{
		{
			scenario: "elem with component",
			elem: &element{
				id:    "hello",
				compo: true,
			},
			byID: "hello",
		},
		{
			scenario: "elem without component",
			elem: &element{
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
			elemDB := NewElementDB()

			if test.elem != nil {
				elemDB.Add(test.elem)

				defer func() {
					elemDB.Remove(test.elem)
					_, err := elemDB.ElementByID(test.elem.ID())
					assert.Error(t, err)
				}()
			}

			elem, err := elemDB.ElementByID(test.byID)
			if test.errID {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.elem.ID(), elem.ID())
			}

			elem, err = elemDB.ElementByComponent(&Bar{})
			if test.errCompo {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.elem.ID(), elem.ID())
			}
		})
	}
}
