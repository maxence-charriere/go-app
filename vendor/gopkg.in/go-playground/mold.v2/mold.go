package mold

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	timeType           = reflect.TypeOf(time.Time{})
	defaultCField      = &cField{}
	restrictedAliasErr = "Alias '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
	restrictedTagErr   = "Tag '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
)

// TODO - ensure StructLevel and Func get passed an interface and not *Transform directly

// Func defines a transform function for use.
type Func func(ctx context.Context, t *Transformer, value reflect.Value, param string) error

// StructLevelFunc accepts all values needed for struct level validation
type StructLevelFunc func(ctx context.Context, t *Transformer, value reflect.Value) error

// Transformer is the base controlling object which contains
// all necessary information
type Transformer struct {
	tagName          string
	aliases          map[string]string
	transformations  map[string]Func
	structLevelFuncs map[reflect.Type]StructLevelFunc
	cCache           *structCache
	tCache           *tagCache
}

// New creates a new Transform object with default tag name of 'mold'
func New() *Transformer {
	tc := new(tagCache)
	tc.m.Store(make(map[string]*cTag))

	sc := new(structCache)
	sc.m.Store(make(map[reflect.Type]*cStruct))

	return &Transformer{
		tagName:         "mold",
		aliases:         make(map[string]string),
		transformations: make(map[string]Func),
		cCache:          sc,
		tCache:          tc,
	}
}

// SetTagName sets the given tag name to be used.
// Default is "trans"
func (t *Transformer) SetTagName(tagName string) {
	t.tagName = tagName
}

// Register adds a transformation with the given tag
//
// NOTES:
// - if the key already exists, the previous transformation function will be replaced.
// - this method is not thread-safe it is intended that these all be registered before hand
func (t *Transformer) Register(tag string, fn Func) {
	if len(tag) == 0 {
		panic("Function Key cannot be empty")
	}

	if fn == nil {
		panic("Function cannot be empty")
	}

	_, ok := restrictedTags[tag]

	if ok || strings.ContainsAny(tag, restrictedTagChars) {
		panic(fmt.Sprintf(restrictedTagErr, tag))
	}

	t.transformations[tag] = fn
}

// RegisterAlias registers a mapping of a single transform tag that
// defines a common or complex set of transformations to simplify adding transforms
// to structs.
//
// NOTE: this function is not thread-safe it is intended that these all be registered before hand
func (t *Transformer) RegisterAlias(alias, tags string) {
	if len(alias) == 0 {
		panic("Alias cannot be empty")
	}

	if len(tags) == 0 {
		panic("Aliased tags cannot be empty")
	}

	_, ok := restrictedTags[alias]

	if ok || strings.ContainsAny(alias, restrictedTagChars) {
		panic(fmt.Sprintf(restrictedAliasErr, alias))
	}
	t.aliases[alias] = tags
}

// RegisterStructLevel registers a StructLevelFunc against a number of types.
// Why does this exist? For structs for which you may not have access or rights to add tags too,
// from other packages your using.
//
// NOTES:
// - this method is not thread-safe it is intended that these all be registered prior to any validation
func (t *Transformer) RegisterStructLevel(fn StructLevelFunc, types ...interface{}) {
	if t.structLevelFuncs == nil {
		t.structLevelFuncs = make(map[reflect.Type]StructLevelFunc)
	}

	for _, typ := range types {
		t.structLevelFuncs[reflect.TypeOf(typ)] = fn
	}
}

// Struct applies transformations against the provided struct
func (t *Transformer) Struct(ctx context.Context, v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr || val.IsNil() {
		return &ErrInvalidTransformValue{typ: reflect.TypeOf(v), fn: "Struct"}
	}

	val = val.Elem()
	typ := val.Type()

	if val.Kind() != reflect.Struct || val.Type() == timeType {
		return &ErrInvalidTransformation{typ: reflect.TypeOf(v)}
	}

	return t.setByStruct(ctx, val, typ, nil)
}

func (t *Transformer) setByStruct(ctx context.Context, current reflect.Value, typ reflect.Type, ct *cTag) (err error) {
	cs, ok := t.cCache.Get(typ)
	if !ok {
		if cs, err = t.extractStructCache(current); err != nil {
			return
		}
	}

	// run is struct has a corresponding struct level transformation
	if cs.fn != nil {
		if err = cs.fn(ctx, t, current); err != nil {
			return
		}
	}

	var f *cField

	for i := 0; i < len(cs.fields); i++ {
		f = cs.fields[i]
		if err = t.setByField(ctx, current.Field(f.idx), f, f.cTags); err != nil {
			return
		}
	}
	return nil
}

// Field applies the provided transformations against the variable
func (t *Transformer) Field(ctx context.Context, v interface{}, tags string) (err error) {
	if len(tags) == 0 || tags == ignoreTag {
		return nil
	}

	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr || val.IsNil() {
		return &ErrInvalidTransformValue{typ: reflect.TypeOf(v), fn: "Field"}
	}
	val = val.Elem()

	// find cached tag
	ctag, ok := t.tCache.Get(tags)
	if !ok {
		t.tCache.lock.Lock()

		// could have been multiple trying to access, but once first is done this ensures tag
		// isn't parsed again.
		ctag, ok = t.tCache.Get(tags)
		if !ok {
			if ctag, _, err = t.parseFieldTagsRecursive(tags, "", "", false); err != nil {
				t.tCache.lock.Unlock()
				return
			}
			t.tCache.Set(tags, ctag)
		}
		t.tCache.lock.Unlock()
	}
	err = t.setByField(ctx, val, defaultCField, ctag)
	return
}

func (t *Transformer) setByField(ctx context.Context, orig reflect.Value, cf *cField, ct *cTag) (err error) {
	current, kind := extractType(orig)

	if ct.hasTag {
		for {
			if ct == nil {
				break
			}

			switch ct.typeof {
			case typeEndKeys:
				return
			case typeDive:
				ct = ct.next

				switch kind {
				case reflect.Slice, reflect.Array:
					reusableCF := &cField{}

					for i := 0; i < current.Len(); i++ {
						if err = t.setByField(ctx, current.Index(i), reusableCF, ct); err != nil {
							return
						}
					}

				case reflect.Map:
					reusableCF := &cField{}

					hasKeys := ct != nil && ct.typeof == typeKeys && ct.keys != nil

					for _, key := range current.MapKeys() {
						newVal := reflect.New(current.Type().Elem()).Elem()
						newVal.Set(current.MapIndex(key))

						if hasKeys {

							// remove current map key as we may be changing it
							// and re-add to the map afterwards
							current.SetMapIndex(key, reflect.Value{})

							newKey := reflect.New(current.Type().Key()).Elem()
							newKey.Set(key)
							key = newKey

							// handle map key
							if err = t.setByField(ctx, key, reusableCF, ct.keys); err != nil {
								return
							}

							// can be nil when just keys being validated
							if ct.next != nil {
								if err = t.setByField(ctx, newVal, reusableCF, ct.next); err != nil {
									return
								}
							}
						} else {
							if err = t.setByField(ctx, newVal, reusableCF, ct); err != nil {
								return
							}
						}
						current.SetMapIndex(key, newVal)
					}

				default:
					err = ErrInvalidDive
				}
				return

			default:
				if !current.CanAddr() {
					newVal := reflect.New(current.Type()).Elem()
					newVal.Set(current)
					if err = ct.fn(ctx, t, newVal, ct.param); err != nil {
						return
					}
					orig.Set(newVal)
				} else {
					if err = ct.fn(ctx, t, current, ct.param); err != nil {
						return
					}
				}
				ct = ct.next
			}
		}
	}

	// need to do this again because one of the previous
	// sets could have set a struct value, where it was a
	// nil pointer before
	current, kind = extractType(current)

	if kind == reflect.Struct {
		typ := current.Type()
		if typ == timeType {
			return
		}
		if ct != nil {
			ct = ct.next
		}

		if !current.CanAddr() {
			newVal := reflect.New(typ).Elem()
			newVal.Set(current)

			if err = t.setByStruct(ctx, newVal, typ, ct); err != nil {
				return
			}
			orig.Set(newVal)
			return
		}
		err = t.setByStruct(ctx, current, typ, ct)
	}
	return
}
