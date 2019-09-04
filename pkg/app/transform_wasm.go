package app

import (
	"strings"
)

type attrTransform func(k, v string) (string, string)

func eventTransform(k, v string) (string, string) {
	if !strings.HasPrefix(k, "on") {
		return k, v
	}
	if strings.HasPrefix(v, "js:") {
		return k, strings.TrimPrefix(v, "js:")
	}
	return k, "//go:" + v
}
