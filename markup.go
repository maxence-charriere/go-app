package app

import (
	"github.com/google/uuid"
)

// Markup is the interface that describes a component set.
// It keeps track of component state and is able to provide info about their
// structure and modifications.
type Markup interface {
	// Component returns the component mounted under the identifier.
	// Returns an error if there is not component with the identifier.
	Component(id uuid.UUID) (c Component, err error)

	// Contains reports whether the component is mounted.
	Contains(c Component) bool

	// Root returns the component root tag.
	// It returns an error if the component is not mounted.
	Root(c Component) (root Tag, err error)

	// Mount indexes the component.
	// The component will be kept in memory until it is dismounted.
	Mount(c Component) (root Tag, err error)

	// Dismount removes references to a component and its children.
	Dismount(c Component)

	// Update updates the tag tree of the component.
	Update(c Component) (syncs []TagSync, err error)
}

// Tag represents a markup tag.
type Tag struct {
	ID       uuid.UUID
	CompoID  uuid.UUID
	Name     string
	Text     string
	Svg      bool
	Type     TagType
	Attrs    AttrMap
	Children []Tag
}

// Is reports whether the tag is of the given type.
func (t *Tag) Is(typ TagType) bool {
	return t.Type == typ
}

// TagType represents a tag type.
type TagType byte

// Constants that enumerates the tag types.
const (
	ZeroTag TagType = iota
	SimpleTag
	TextTag
	CompoTag
)

// AttrMap represents a map of attributes.
type AttrMap map[string]string

// TagSync represents a tag synchronisation.
type TagSync struct {
	Tag     *Tag
	Replace bool
}
