package objconv

import (
	"reflect"
	"sync"
)

// An Adapter is a pair of an encoder and a decoder function that can be
// installed on the package to support new types.
type Adapter struct {
	Encode func(Encoder, reflect.Value) error
	Decode func(Decoder, reflect.Value) error
}

// Install adds an adapter for typ.
//
// The function panics if one of the encoder and decoder functions of the
// adapter are nil.
//
// A typical use case for this function is to be called during the package
// initialization phase to extend objconv support for new types.
func Install(typ reflect.Type, adapter Adapter) {
	if adapter.Encode == nil {
		panic("objconv: the encoder function of an adapter cannot be nil")
	}

	if adapter.Decode == nil {
		panic("objconv: the decoder function of an adapter cannot be nil")
	}

	adapterMutex.Lock()
	adapterStore[typ] = adapter
	adapterMutex.Unlock()

	// We have to clear the struct cache because it may now have become invalid.
	// Because installing adapters is done in the package initialization phase
	// it's unlikely that any encoding or decoding operations are taking place
	// at this time so there should be no performance impact of clearing the
	// cache.
	structCache.clear()
}

// AdapterOf returns the adapter for typ, setting ok to true if one was found,
// false otherwise.
func AdapterOf(typ reflect.Type) (a Adapter, ok bool) {
	adapterMutex.RLock()
	a, ok = adapterStore[typ]
	adapterMutex.RUnlock()
	return
}

var (
	adapterMutex sync.RWMutex
	adapterStore = make(map[reflect.Type]Adapter)
)
