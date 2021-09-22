package cache

import (
	"context"
	"reflect"
)

// Cache is the interface that describes a cache.
type Cache interface {
	// Get returns the item with the given key, otherwise returns false.
	Get(ctx context.Context, key string) (Item, bool)

	// Set sets the item at the given key.
	Set(ctx context.Context, key string, i Item)

	// Deletes the item at the given key.
	Del(ctx context.Context, key string)

	// The number of items in the cache.
	Len() int

	// The size in bytes.
	Size() int
}

// Item is the interface that describes a cacheable item.
type Item interface {
	// The size that the item occupies in a cache.
	Size() int
}

// Bytes represents a cacheable byte slice.
type Bytes []byte

func (b Bytes) Size() int {
	return len(b)
}

type String string

func (s String) Size() int {
	return len(s)
}

type Int int

func (i Int) Size() int {
	return intSize
}

type Float float64

func (f Float) Size() int {
	return floatSize
}

var (
	intSize   = int(reflect.TypeOf(42).Size())
	floatSize = int(reflect.TypeOf(23.42).Size())
)
