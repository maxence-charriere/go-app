package net

import (
	"net"
	"reflect"

	"github.com/segmentio/objconv"
)

func init() {
	objconv.Install(reflect.TypeOf(net.TCPAddr{}), TCPAddrAdapter())
	objconv.Install(reflect.TypeOf(net.UDPAddr{}), UDPAddrAdapter())
	objconv.Install(reflect.TypeOf(net.UnixAddr{}), UnixAddrAdapter())
	objconv.Install(reflect.TypeOf(net.IPAddr{}), IPAddrAdapter())
	objconv.Install(reflect.TypeOf(net.IP(nil)), IPAdapter())
}

// TCPAddrAdapter returns the adapter to encode and decode net.TCPAddr values.
func TCPAddrAdapter() objconv.Adapter {
	return objconv.Adapter{
		Encode: encodeTCPAddr,
		Decode: decodeTCPAddr,
	}
}

// UDPAddrAdapter returns the adapter to encode and decode net.UDPAddr values.
func UDPAddrAdapter() objconv.Adapter {
	return objconv.Adapter{
		Encode: encodeUDPAddr,
		Decode: decodeUDPAddr,
	}
}

// UnixAddrAdapter returns the adapter to encode and decode net.UnixAddr values.
func UnixAddrAdapter() objconv.Adapter {
	return objconv.Adapter{
		Encode: encodeUnixAddr,
		Decode: decodeUnixAddr,
	}
}

// IPAddrAdapter returns the adapter to encode and decode net.IPAddr values.
func IPAddrAdapter() objconv.Adapter {
	return objconv.Adapter{
		Encode: encodeIPAddr,
		Decode: decodeIPAddr,
	}
}

// IPAdapter returns the adapter to encode and decode net.IP values.
func IPAdapter() objconv.Adapter {
	return objconv.Adapter{
		Encode: encodeIP,
		Decode: decodeIP,
	}
}
