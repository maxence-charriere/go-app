package app

// Tagger is the interface that describes a collection of tags that gives
// context to something.
type Tagger interface {
	// Returns a collection of tags.
	Tags() Tags
}

// Tags represent key-value pairs that give context to what they are used with.
type Tags map[string]string

func (t Tags) Tags() Tags {
	return t
}

// Set sets a tag with the given name and value. The value is converted to a
// string.
func (t Tags) Set(name string, v any) {
	t[name] = toString(v)
}

// Get returns a tag value with the given name.
func (t Tags) Get(name string) string {
	return t[name]
}

// Tag is a key-value pair that adds context to an action.
type Tag struct {
	Name  string
	Value string
}

func (t Tag) Tags() Tags {
	return Tags{t.Name: t.Value}
}

// T creates a tag with the given name and value. The value is converted to a
// string.
func T(name string, value any) Tag {
	return Tag{
		Name:  name,
		Value: toString(value),
	}
}
