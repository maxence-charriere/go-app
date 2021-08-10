package app

// Tags represent key-value pairs that give context to what they are used with.
type Tags map[string]string

func makeTags(tags []Tag) Tags {
	if len(tags) == 0 {
		return nil
	}

	t := make(Tags, len(tags))
	for _, tag := range tags {
		t[tag.Name] = tag.Value
	}
	return t
}

// Set sets a tag with the given name and value. The value is converted to a
// string.
func (t Tags) Set(name string, v interface{}) {
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

// T creates a tag with the given name and value. The value is converted to a
// string.
func T(name string, value interface{}) Tag {
	return Tag{
		Name:  name,
		Value: toString(value),
	}
}
