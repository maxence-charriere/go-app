package url

import (
	"errors"
	"net/url"
	"reflect"

	"github.com/segmentio/objconv"
)

func decodeURL(d objconv.Decoder, to reflect.Value) (err error) {
	var u *url.URL
	var s string

	if err = d.Decode(&s); err != nil {
		return
	}

	if u, err = url.Parse(s); err != nil {
		err = errors.New("objconv: bad URL: " + err.Error())
		return
	}

	if to.IsValid() {
		to.Set(reflect.ValueOf(*u))
	}
	return
}

func decodeQuery(d objconv.Decoder, to reflect.Value) (err error) {
	var v url.Values
	var s string

	if err = d.Decode(&s); err != nil {

	}

	if v, err = url.ParseQuery(s); err != nil {
		err = errors.New("objconv: bad URL values: " + err.Error())
		return
	}

	if to.IsValid() {
		to.Set(reflect.ValueOf(v))
	}
	return
}
