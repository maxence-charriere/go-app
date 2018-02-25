package mold

import (
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

type tagType uint8

const (
	typeDefault tagType = iota
	typeDive
	typeKeys
	typeEndKeys
)

const (
	keysTagNotDefined = "'" + endKeysTag + "' tag encountered without a corresponding '" + keysTag + "' tag"
)

type structCache struct {
	lock sync.Mutex
	m    atomic.Value // map[reflect.Type]*cStruct
}

func (sc *structCache) Get(key reflect.Type) (c *cStruct, found bool) {
	c, found = sc.m.Load().(map[reflect.Type]*cStruct)[key]
	return
}

func (sc *structCache) Set(key reflect.Type, value *cStruct) {

	m := sc.m.Load().(map[reflect.Type]*cStruct)

	nm := make(map[reflect.Type]*cStruct, len(m)+1)
	for k, v := range m {
		nm[k] = v
	}
	nm[key] = value
	sc.m.Store(nm)
}

type tagCache struct {
	lock sync.Mutex
	m    atomic.Value // map[string]*cTag
}

func (tc *tagCache) Get(key string) (c *cTag, found bool) {
	c, found = tc.m.Load().(map[string]*cTag)[key]
	return
}

func (tc *tagCache) Set(key string, value *cTag) {

	m := tc.m.Load().(map[string]*cTag)

	nm := make(map[string]*cTag, len(m)+1)
	for k, v := range m {
		nm[k] = v
	}
	nm[key] = value
	tc.m.Store(nm)
}

type cStruct struct {
	fields []*cField
	fn     StructLevelFunc
}

type cField struct {
	idx   int
	cTags *cTag
}

type cTag struct {
	tag            string
	aliasTag       string
	actualAliasTag string
	hasAlias       bool
	typeof         tagType
	hasTag         bool
	fn             Func
	keys           *cTag
	next           *cTag
	param          string
}

func (t *Transformer) extractStructCache(current reflect.Value) (*cStruct, error) {
	t.cCache.lock.Lock()
	defer t.cCache.lock.Unlock()

	typ := current.Type()

	// could have been multiple trying to access, but once first is done this ensures struct
	// isn't parsed again.
	cs, ok := t.cCache.Get(typ)
	if ok {
		return cs, nil
	}

	cs = &cStruct{fields: make([]*cField, 0), fn: t.structLevelFuncs[typ]}
	numFields := current.NumField()

	var ctag *cTag
	var fld reflect.StructField
	var tag string
	var err error

	for i := 0; i < numFields; i++ {

		fld = typ.Field(i)

		if !fld.Anonymous && len(fld.PkgPath) > 0 {
			continue
		}

		tag = fld.Tag.Get(t.tagName)
		if tag == ignoreTag {
			continue
		}

		// NOTE: cannot use shared tag cache, because tags may be equal, but things like alias may be different
		// and so only struct level caching can be used instead of combined with Field tag caching
		if len(tag) > 0 {
			ctag, _, err = t.parseFieldTagsRecursive(tag, fld.Name, "", false)
			if err != nil {
				return nil, err
			}
		} else {
			// even if field doesn't have validations need cTag for traversing to potential inner/nested
			// elements of the field.
			ctag = new(cTag)
		}

		cs.fields = append(cs.fields, &cField{
			idx:   i,
			cTags: ctag,
		})
	}

	t.cCache.Set(typ, cs)

	return cs, nil
}

func (t *Transformer) parseFieldTagsRecursive(tag string, fieldName string, alias string, hasAlias bool) (firstCtag *cTag, current *cTag, err error) {

	var tg string
	var ok bool
	noAlias := len(alias) == 0
	tags := strings.Split(tag, tagSeparator)

	for i := 0; i < len(tags); i++ {

		tg = tags[i]
		if noAlias {
			alias = tg
		}

		// check map for alias and process new tags, otherwise process as usual
		if tagsVal, found := t.aliases[tg]; found {
			if i == 0 {
				firstCtag, current, err = t.parseFieldTagsRecursive(tagsVal, fieldName, tg, true)
				if err != nil {
					return
				}
			} else {
				next, curr, errr := t.parseFieldTagsRecursive(tagsVal, fieldName, tg, true)
				if errr != nil {
					err = errr
					return
				}
				current.next, current = next, curr
			}
			continue
		}

		var prevTag tagType

		if i == 0 {
			current = &cTag{aliasTag: alias, hasAlias: hasAlias, hasTag: true}
			firstCtag = current
		} else {
			prevTag = current.typeof
			current.next = &cTag{aliasTag: alias, hasAlias: hasAlias, hasTag: true}
			current = current.next
		}

		switch tg {

		case diveTag:
			current.typeof = typeDive
			continue

		case keysTag:
			current.typeof = typeKeys

			if i == 0 || prevTag != typeDive {
				err = ErrInvalidKeysTag
				return
			}

			current.typeof = typeKeys

			// need to pass along only keys tag
			// need to increment i to skip over the keys tags
			b := make([]byte, 0, 64)

			i++

			for ; i < len(tags); i++ {

				b = append(b, tags[i]...)
				b = append(b, ',')

				if tags[i] == endKeysTag {
					break
				}
			}

			if current.keys, _, err = t.parseFieldTagsRecursive(string(b[:len(b)-1]), fieldName, "", false); err != nil {
				return
			}
			continue

		case endKeysTag:
			current.typeof = typeEndKeys

			// if there are more in tags then there was no keysTag defined
			// and an error should be thrown
			if i != len(tags)-1 {
				err = ErrUndefinedKeysTag
			}
			return

		default:

			vals := strings.SplitN(tg, tagKeySeparator, 2)

			if noAlias {
				alias = vals[0]
				current.aliasTag = alias
			} else {
				current.actualAliasTag = tg
			}

			current.tag = vals[0]
			if len(current.tag) == 0 {
				err = &ErrInvalidTag{tag: current.tag, field: fieldName}
				return
			}

			if current.fn, ok = t.transformations[current.tag]; !ok {
				err = &ErrUndefinedTag{tag: current.tag, field: fieldName}
				return
			}

			if len(vals) > 1 {
				current.param = strings.Replace(vals[1], utf8HexComma, ",", -1)
			}
		}
	}
	return
}
