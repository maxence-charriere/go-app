package mail

import (
	"bytes"
	"net/mail"
	"reflect"

	"github.com/segmentio/objconv"
)

func encodeAddress(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(mail.Address)
	return e.Encode(a.String())
}

func encodeAddressList(e objconv.Encoder, v reflect.Value) error {
	l := v.Interface().([]*mail.Address)
	b := &bytes.Buffer{}

	for i, a := range l {
		if a == nil {
			continue
		}
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(a.String())
	}

	return e.Encode(b.String())
}
