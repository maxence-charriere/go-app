package html

import (
	"encoding/json"
	"reflect"
	"strings"
	"unicode"

	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

type mapper struct {
	completePipeline []string
	index            int
	jsonValue        string
}

func newMapper(pipeline []string, jsonValue string) *mapper {
	return &mapper{
		completePipeline: pipeline,
		jsonValue:        jsonValue,
	}
}

func (m *mapper) pipeline() []string {
	return m.completePipeline[m.index:]
}

func (m *mapper) target() string {
	return m.completePipeline[m.index]
}

func (m *mapper) fullTarget() string {
	return strings.Join(m.completePipeline[:m.index+1], ".")
}

func (m *mapper) MapTo(compo app.Component) (funcMapping bool, err error) {
	return m.mapTo(reflect.ValueOf(compo))
}

func (m *mapper) mapTo(value reflect.Value) (funcMapping bool, err error) {
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

func (m *mapper) mapToPointer(ptr reflect.Value) (funcMapping bool, err error) {
	if len(m.pipeline()) == 0 {
		return m.mapToValue(ptr)
	}

	if !isExported(m.target()) {
		err = errors.Errorf("%s is mapped to an unexported method", m.fullTarget())
		return
	}

	method := ptr.MethodByName(m.target())
	if method.IsValid() {
		m.index++
		return m.mapTo(method)
	}

	return m.mapTo(ptr.Elem())
}

func (m *mapper) mapToStruct(structure reflect.Value) (funcMapping bool, err error) {
	if len(m.pipeline()) == 0 {
		return m.mapToValue(structure)
	}

	if !isExported(m.target()) {
		err = errors.Errorf(
			"%s is mapped to unexported field or method",
			m.fullTarget(),
		)
		return
	}

	if method := structure.MethodByName(m.target()); method.IsValid() {
		m.index++
		return m.mapTo(method)
	}

	field := structure.FieldByName(m.target())
	if !field.IsValid() {
		err = errors.Errorf(
			"%s is mapped to a nonexistent field or method",
			m.fullTarget(),
		)
		return
	}

	m.index++
	return m.mapTo(field)
}

func (m *mapper) mapToMap(mapv reflect.Value) (funcMapping bool, err error) {
	if len(m.pipeline()) == 0 {
		m.mapToValue(mapv)
		return
	}

	if isExported(m.target()) {
		if method := mapv.MethodByName(m.target()); method.IsValid() {
			m.index++
			return m.mapTo(method)
		}
	}

	err = errors.Errorf(
		"%s is mapped to a map value",
		m.fullTarget(),
	)
	return
}

func (m *mapper) mapToSlice(slice reflect.Value) (funcMapping bool, err error) {
	if len(m.pipeline()) == 0 {
		return m.mapToValue(slice)
	}

	if isExported(m.target()) {
		if method := slice.MethodByName(m.target()); method.IsValid() {
			m.index++
			return m.mapTo(method)
		}
	}

	err = errors.Errorf(
		"%s is mapped to a slice value",
		m.fullTarget(),
	)
	return
}

func (m *mapper) mapToFunction(function reflect.Value) (funcMapping bool, err error) {
	if len(m.pipeline()) != 0 {
		err = errors.Errorf(
			"%s is mapped to a unsuported method",
			m.fullTarget(),
		)
		return
	}

	typ := function.Type()
	if typ.NumIn() > 1 {
		err = errors.Errorf(
			"%s is mapped to func that have more than 1 arg",
			m.completePipeline,
		)
		return
	}

	funcMapping = true

	if typ.NumIn() == 0 {
		function.Call(nil)
		return
	}

	arg := reflect.New(typ.In(0))

	if err = json.Unmarshal([]byte(m.jsonValue), arg.Interface()); err != nil {
		err = errors.Wrapf(err, "%s:", m.completePipeline)
		return
	}

	function.Call([]reflect.Value{arg.Elem()})
	return
}

func (m *mapper) mapToValue(value reflect.Value) (funcMapping bool, err error) {
	if len(m.pipeline()) == 0 {
		newValue := reflect.New(value.Type())

		if err = json.Unmarshal([]byte(m.jsonValue), newValue.Interface()); err != nil {
			err = errors.Wrapf(err, "%s:", m.completePipeline)
			return
		}

		value.Set(newValue.Elem())
		return
	}

	if !isExported(m.target()) {
		err = errors.Errorf(
			"%s is mapped to a unsuported method",
			m.fullTarget(),
		)
		return
	}

	method := value.MethodByName(m.target())
	if !method.IsValid() {
		err = errors.Errorf(
			"%s is mapped to a undefined method",
			m.fullTarget(),
		)
		return
	}

	m.index++
	return m.mapTo(method)
}

func isExported(fieldOrMethod string) bool {
	return !unicode.IsLower(rune(fieldOrMethod[0]))
}
