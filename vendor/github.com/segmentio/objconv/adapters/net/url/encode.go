package url

import (
	"net/url"
	"reflect"

	"github.com/segmentio/objconv"
)

func encodeURL(e objconv.Encoder, v reflect.Value) error {
	u := v.Interface().(url.URL)
	return e.Encode(u.String())
}

func encodeQuery(e objconv.Encoder, v reflect.Value) error {
	q := v.Interface().(url.Values)
	return e.Encode(q.Encode())
}
