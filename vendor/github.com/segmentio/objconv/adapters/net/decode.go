package net

import (
	"errors"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/segmentio/objconv"
)

func decodeTCPAddr(d objconv.Decoder, to reflect.Value) (err error) {
	var a net.TCPAddr
	var s string

	if err = d.Decode(&s); err != nil {
		return
	}

	if a.IP, a.Port, a.Zone, err = parseNetAddr(s); err != nil {
		return
	}

	if to.IsValid() {
		to.Set(reflect.ValueOf(a))
	}
	return
}

func decodeUDPAddr(d objconv.Decoder, to reflect.Value) (err error) {
	var a net.UDPAddr
	var s string

	if err = d.Decode(&s); err != nil {
		return
	}

	if a.IP, a.Port, a.Zone, err = parseNetAddr(s); err != nil {
		return
	}

	if to.IsValid() {
		to.Set(reflect.ValueOf(a))
	}
	return
}

func decodeUnixAddr(d objconv.Decoder, to reflect.Value) (err error) {
	var a net.UnixAddr
	var s string

	if err = d.Decode(&s); err != nil {
		return
	}

	if i := strings.Index(s, "://"); i >= 0 {
		a.Net, a.Name = s[:i], s[i+3:]
	} else {
		a.Net, a.Name = "unix", s
	}

	if to.IsValid() {
		to.Set(reflect.ValueOf(a))
	}
	return
}

func decodeIPAddr(d objconv.Decoder, to reflect.Value) (err error) {
	var a net.IPAddr
	var s string

	if err = d.Decode(&s); err != nil {
		return
	}

	if i := strings.IndexByte(s, '%'); i >= 0 {
		s, a.Zone = s[:i], s[i+1:]
	}

	if a.IP = net.ParseIP(s); a.IP == nil {
		err = errors.New("objconv: bad IP address: " + s)
		return
	}

	if to.IsValid() {
		to.Set(reflect.ValueOf(a))
	}
	return
}

func decodeIP(d objconv.Decoder, to reflect.Value) (err error) {
	var ip net.IP
	var s string

	if err = d.Decode(&s); err != nil {
		return
	}

	if ip = net.ParseIP(s); ip == nil {
		err = errors.New("objconv: bad IP address: " + s)
		return
	}

	if to.IsValid() {
		to.Set(reflect.ValueOf(ip))
	}
	return
}

func parseNetAddr(s string) (ip net.IP, port int, zone string, err error) {
	var h string
	var p string

	if h, p, err = net.SplitHostPort(s); err != nil {
		h, p = s, ""
	}

	if len(h) != 0 {
		if off := strings.IndexByte(h, '%'); off >= 0 {
			h, zone = h[:off], h[off+1:]
		}
		if ip = net.ParseIP(h); ip == nil {
			err = errors.New("objconv: bad IP address: " + s)
			return
		}
	}

	if len(p) != 0 {
		if port, err = strconv.Atoi(p); err != nil || port < 0 || port > 65535 {
			err = errors.New("objconv: bad port number: " + s)
			return
		}
	}

	return
}
