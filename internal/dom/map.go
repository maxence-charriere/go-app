package dom

import (
	"encoding/json"
	"reflect"
	"strings"
	"unicode"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

// Mapping represents a component method or field descriptor.
type Mapping struct {
	// The component identifier.
	CompoID string

	// A dot separated string that points to a component field or method.
	FieldOrMethod string

	// The JSON value to map to a field or method's first argument.
	JSONValue string

	// A string that describes a field that may required override.
	Override string

	pipeline []string
	index    int
}

// Map performs the mapping to the given component.
func (m *Mapping) Map(c app.Compo) (f func(), err error) {
	if m.pipeline, err = pipeline(m.FieldOrMethod); err != nil {
		return nil, err
	}

	return m.mapTo(reflect.ValueOf(c))
}

func (m *Mapping) currentPipeline() []string {
	return m.pipeline[m.index:]
}

func (m *Mapping) target() string {
	return m.pipeline[m.index]
}

func (m *Mapping) fullTarget() string {
	return strings.Join(m.pipeline[:m.index+1], ".")
}

func (m *Mapping) mapTo(value reflect.Value) (func(), error) {
	switch value.Kind() {
	case reflect.Ptr:
		return m.mapToPointer(value)

	case reflect.Struct:
		return m.mapToStruct(value)

	case reflect.Map:
		return m.mapToMap(value)

	case reflect.Slice, reflect.Array:
		return m.mapToSlice(value)

	case reflect.Func:
		return m.mapToFunction(value)

	default:
		return m.mapToValue(value)
	}
}

func (m *Mapping) mapToPointer(ptr reflect.Value) (func(), error) {
	if len(m.currentPipeline()) == 0 {
		return m.mapToValue(ptr)
	}

	if !isExported(m.target()) {
		return nil, errors.Errorf("%s is mapped to an unexported method", m.fullTarget())
	}

	method := ptr.MethodByName(m.target())
	if method.IsValid() {
		m.index++
		return m.mapTo(method)
	}

	return m.mapTo(ptr.Elem())
}

func (m *Mapping) mapToStruct(structure reflect.Value) (func(), error) {
	if len(m.currentPipeline()) == 0 {
		return m.mapToValue(structure)
	}

	if !isExported(m.target()) {
		return nil, errors.Errorf(
			"%s is mapped to unexported field or method",
			m.fullTarget(),
		)
	}

	if method := structure.MethodByName(m.target()); method.IsValid() {
		m.index++
		return m.mapTo(method)
	}

	field := structure.FieldByName(m.target())
	if !field.IsValid() {
		return nil, errors.Errorf(
			"%s is mapped to a nonexistent field or method",
			m.fullTarget(),
		)
	}

	m.index++
	return m.mapTo(field)
}

func (m *Mapping) mapToMap(mapv reflect.Value) (func(), error) {
	if len(m.currentPipeline()) == 0 {
		return m.mapToValue(mapv)
	}

	if isExported(m.target()) {
		if method := mapv.MethodByName(m.target()); method.IsValid() {
			m.index++
			return m.mapTo(method)
		}
	}

	return nil, errors.Errorf(
		"%s is mapped to a map value",
		m.fullTarget(),
	)
}

func (m *Mapping) mapToSlice(slice reflect.Value) (func(), error) {
	if len(m.currentPipeline()) == 0 {
		return m.mapToValue(slice)
	}

	if isExported(m.target()) {
		if method := slice.MethodByName(m.target()); method.IsValid() {
			m.index++
			return m.mapTo(method)
		}
	}

	return nil, errors.Errorf(
		"%s is mapped to a slice value",
		m.fullTarget(),
	)
}

func (m *Mapping) mapToFunction(fn reflect.Value) (func(), error) {
	if len(m.currentPipeline()) != 0 {
		return nil, errors.Errorf(
			"%s is mapped to a unsuported method",
			m.fullTarget(),
		)
	}

	typ := fn.Type()
	if typ.NumIn() > 1 {
		return nil, errors.Errorf(
			"%s is mapped to func that have more than 1 arg",
			m.pipeline,
		)
	}

	if typ.NumIn() == 0 {
		return func() {
			fn.Call(nil)
		}, nil
	}

	arg := reflect.New(typ.In(0))

	if err := json.Unmarshal([]byte(m.JSONValue), arg.Interface()); err != nil {
		return nil, errors.Wrapf(err, "%s:", m.pipeline)
	}

	return func() {
		fn.Call([]reflect.Value{arg.Elem()})
	}, nil
}

func (m *Mapping) mapToValue(value reflect.Value) (func(), error) {
	if len(m.currentPipeline()) == 0 {
		newValue := reflect.New(value.Type())

		if err := json.Unmarshal([]byte(m.JSONValue), newValue.Interface()); err != nil {
			return nil, errors.Wrapf(err, "%s:", m.pipeline)
		}

		value.Set(newValue.Elem())
		return nil, nil
	}

	if !isExported(m.target()) {
		return nil, errors.Errorf(
			"%s is mapped to a unsuported method",
			m.fullTarget(),
		)
	}

	method := value.MethodByName(m.target())
	if !method.IsValid() {
		return nil, errors.Errorf(
			"%s is mapped to a undefined method",
			m.fullTarget(),
		)
	}

	m.index++
	return m.mapTo(method)
}

func pipeline(fieldOrMethod string) ([]string, error) {
	if len(fieldOrMethod) == 0 {
		return nil, errors.New("empty")
	}

	p := strings.Split(fieldOrMethod, ".")

	for _, e := range p {
		if len(e) == 0 {
			return nil, errors.Errorf("%s: contains an empty element", fieldOrMethod)
		}
	}

	return p, nil
}

func isExported(fieldOrMethod string) bool {
	return !unicode.IsLower(rune(fieldOrMethod[0]))
}
