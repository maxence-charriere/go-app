package app

// Tags represent key-value pairs that give context to what they are used with.
type Tags map[string]string

// Set sets a tag with the given name and value. The value is converted to a
// string.
func (t Tags) Set(name string, v interface{}) {
	t[name] = toString(v)
}

// Get returns a tag value with the given name.
func (t Tags) Get(name string) string {
	return t[name]
}
