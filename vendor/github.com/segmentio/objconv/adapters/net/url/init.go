package url

import (
	"net/url"
	"reflect"

	"github.com/segmentio/objconv"
)

func init() {
	objconv.Install(reflect.TypeOf(url.URL{}), URLAdapter())
	objconv.Install(reflect.TypeOf(url.Values(nil)), QueryAdapter())
}

// URLAdapter returns the adapter to encode and decode url.URL values.
func URLAdapter() objconv.Adapter {
	return objconv.Adapter{
		Encode: encodeURL,
		Decode: decodeURL,
	}
}

// QueryAdapter returns the adapter to encode and decode url.Values values.
func QueryAdapter() objconv.Adapter {
	return objconv.Adapter{
		Encode: encodeQuery,
		Decode: decodeQuery,
	}
}
