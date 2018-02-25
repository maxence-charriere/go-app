package objconv

import (
	"reflect"
	"sync"

	"github.com/segmentio/objconv/objutil"
)

// structField represents a single field of a struct and carries information
// useful to the algorithms of the objconv package.
type structField struct {
	// The index of the field in the structure.
	index []int

	// The name of the field in the structure.
	name string

	// Omitempty is set to true when the field should be omitted if it has an
	// empty value.
	omitempty bool

	// Omitzero is set to true when the field should be omitted if it has a zero
	// value.
	omitzero bool

	// cache for the encoder and decoder methods
	encode encodeFunc
	decode decodeFunc
}

func makeStructField(f reflect.StructField, c map[reflect.Type]*structType) structField {
	var t objutil.Tag

	if tag := f.Tag.Get("objconv"); len(tag) != 0 {
		t = objutil.ParseTag(tag)
	} else {
		// To maximize compatibility with existing code we fallback to checking
		// if the field has a `json` tag.
		//
		// This tag doesn't support any of the extra features that are supported
		// by the `objconv` tag, and it should stay this way. It has to match
		// the behavior of the standard encoding/json package to avoid any
		// implicit changes in what would be intuitively expected.
		t = objutil.ParseTagJSON(f.Tag.Get("json"))
	}

	s := structField{
		index:     f.Index,
		name:      f.Name,
		omitempty: t.Omitempty,
		omitzero:  t.Omitzero,

		encode: makeEncodeFunc(f.Type, encodeFuncOpts{
			recurse: true,
			structs: c,
		}),

		decode: makeDecodeFunc(f.Type, decodeFuncOpts{
			recurse: true,
			structs: c,
		}),
	}

	if len(t.Name) != 0 {
		s.name = t.Name
	}

	return s
}

func (f *structField) omit(v reflect.Value) bool {
	return (f.omitempty && objutil.IsEmptyValue(v)) || (f.omitzero && objutil.IsZeroValue(v))
}

// structType is used to represent a Go structure in internal data structures
// that cache meta information to make field lookups faster and avoid having to
// use reflection to lookup the same type information over and over again.
type structType struct {
	fields       []structField           // the serializable fields of the struct
	fieldsByName map[string]*structField // cache of fields by name
}

// newStructType takes a Go type as argument and extract information to make a
// new structType value.
// The type has to be a struct type or a panic will be raised.
func newStructType(t reflect.Type, c map[reflect.Type]*structType) *structType {
	if s := c[t]; s != nil {
		return s
	}

	n := t.NumField()
	s := &structType{
		fields:       make([]structField, 0, n),
		fieldsByName: make(map[string]*structField),
	}
	c[t] = s

	for i := 0; i != n; i++ {
		ft := t.Field(i)

		if ft.Anonymous || len(ft.PkgPath) != 0 { // anonymous or non-exported
			continue
		}

		sf := makeStructField(ft, c)

		if sf.name == "-" { // skip
			continue
		}

		s.fields = append(s.fields, sf)
		s.fieldsByName[sf.name] = &s.fields[len(s.fields)-1]
	}

	return s
}

// structTypeCache is a simple cache for mapping Go types to Struct values.
type structTypeCache struct {
	mutex sync.RWMutex
	store map[reflect.Type]*structType
}

// lookup takes a Go type as argument and returns the matching structType value,
// potentially creating it if it didn't already exist.
// This method is safe to call from multiple goroutines.
func (cache *structTypeCache) lookup(t reflect.Type) (s *structType) {
	cache.mutex.RLock()
	s = cache.store[t]
	cache.mutex.RUnlock()

	if s == nil {
		// There's a race confition here where this value may be generated
		// multiple times.
		// The impact in practice is really small as it's unlikely to happen
		// often, we take the approach of keeping the logic simple and avoid
		// a more complex synchronization logic required to solve this edge
		// case.
		s = newStructType(t, map[reflect.Type]*structType{})
		cache.mutex.Lock()
		cache.store[t] = s
		cache.mutex.Unlock()
	}

	return
}

// clear empties the cache.
func (cache *structTypeCache) clear() {
	cache.mutex.Lock()
	for typ := range cache.store {
		delete(cache.store, typ)
	}
	cache.mutex.Unlock()
}

var (
	// This struct cache is used to avoid reusing reflection over and over when
	// the objconv functions are called. The performance improvements on iterating
	// over struct fields are huge, this is a really important optimization:
	structCache = structTypeCache{
		store: make(map[reflect.Type]*structType),
	}
)
