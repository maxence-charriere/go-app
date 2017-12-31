package app

import (
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Markup is the interface that describes a component set.
// It keeps track of component states and is able to provide info about their
// structure and modifications.
type Markup interface {
	// Len returns the number of components living in the markup.
	Len() int

	// Component returns the component mounted under the identifier.
	// Returns an error if there is not component with the identifier.
	Component(id uuid.UUID) (compo Component, err error)

	// Contains reports whether the component is mounted.
	Contains(compo Component) bool

	// Root returns the component root tag.
	// It returns an error if the component is not mounted.
	Root(compo Component) (root Tag, err error)

	// Mount indexes the component.
	// The component will be kept in memory until it is dismounted.
	Mount(compo Component) (root Tag, err error)

	// Dismount removes references to a component and its children.
	Dismount(compo Component)

	// Update updates the tag tree of the component.
	Update(compo Component) (syncs []TagSync, err error)

	// Map performs the given mapping.
	// The json value is mapped to the field or method of the specified
	// component.
	// Methods and fields of func type are called with the value mapped to their
	// first arg.
	// It returns an error if the assigned field or method is not exported.
	Map(mapping Mapping) (function func(), err error)
}

// Mapping describes a component mapping.
type Mapping struct {
	// The component identifier.
	CompoID uuid.UUID

	// A dot separated string that points to a component field or method.
	Target string

	// The JSON value to map to a field or method's first argument.
	JSONValue string
}

// ParseMappingTarget parses the given target and returns the corresponding
// pipeline.
func ParseMappingTarget(target string) (pipeline []string, err error) {
	if len(target) == 0 {
		err = errors.New("empty target")
	}

	pipeline = strings.Split(target, ".")

	for _, elem := range pipeline {
		if len(elem) == 0 {
			err = errors.Errorf("%s contains empty element", target)
			return
		}
	}
	return
}

// TagEncoder is the interface that describes an encoder that writes the tag
// markup representation to an output stream.
type TagEncoder interface {
	// Encode write the tag as a markup representation to its output.
	Encode(tag Tag) error
}

// TagDecoder is the interface that describes a decoder that reads and decodes
// tags from an input stream.
type TagDecoder interface {
	// Decode reads the markup from its input put and store it in the given tag.
	Decode(tag *Tag) error
}

// Tag represents a markup tag.
type Tag struct {
	ID         uuid.UUID
	CompoID    uuid.UUID
	Name       string
	Text       string
	Svg        bool
	Type       TagType
	Attributes AttributeMap
	Children   []Tag
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

// AttributeMap represents a map of attributes.
type AttributeMap map[string]string

// TagSync represents a tag synchronisation.
type TagSync struct {
	Tag     Tag
	Replace bool
}

// NewConcurrentMarkup decorates the given markup to ensure concurrent access
// safety.
func NewConcurrentMarkup(markup Markup) Markup {
	return &concurrentMarkup{
		base: markup,
	}
}

type concurrentMarkup struct {
	mutex sync.Mutex
	base  Markup
}

func (m *concurrentMarkup) Len() int {
	m.mutex.Lock()
	l := m.base.Len()
	m.mutex.Unlock()
	return l
}

func (m *concurrentMarkup) Component(id uuid.UUID) (compo Component, err error) {
	m.mutex.Lock()
	compo, err = m.base.Component(id)
	m.mutex.Unlock()
	return
}

func (m *concurrentMarkup) Contains(compo Component) bool {
	m.mutex.Lock()
	contains := m.base.Contains(compo)
	m.mutex.Unlock()
	return contains
}

func (m *concurrentMarkup) Root(compo Component) (root Tag, err error) {
	m.mutex.Lock()
	root, err = m.base.Root(compo)
	m.mutex.Unlock()
	return
}

func (m *concurrentMarkup) Mount(compo Component) (root Tag, err error) {
	m.mutex.Lock()
	root, err = m.base.Mount(compo)
	m.mutex.Unlock()
	return
}

func (m *concurrentMarkup) Dismount(compo Component) {
	m.mutex.Lock()
	m.base.Dismount(compo)
	m.mutex.Unlock()
}

func (m *concurrentMarkup) Update(compo Component) (syncs []TagSync, err error) {
	m.mutex.Lock()
	syncs, err = m.base.Update(compo)
	m.mutex.Unlock()
	return
}

func (m *concurrentMarkup) Map(mapping Mapping) (function func(), err error) {
	m.mutex.Lock()
	function, err = m.base.Map(mapping)
	m.mutex.Unlock()
	return
}
