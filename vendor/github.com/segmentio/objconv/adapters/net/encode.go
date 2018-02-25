package net

import (
	"net"
	"reflect"

	"github.com/segmentio/objconv"
)

func encodeTCPAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.TCPAddr)
	return e.Encode(a.String())
}

func encodeUDPAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.UDPAddr)
	return e.Encode(a.String())
}

func encodeUnixAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.UnixAddr)
	return e.Encode(a.String())
}

func encodeIPAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.IPAddr)
	return e.Encode(a.String())
}

func encodeIP(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.IP)
	return e.Encode(a.String())
}
