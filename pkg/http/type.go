package http

import "strings"

// DefaultContentTypes contains the content types that are cacheable or
// compressible by default.
var DefaultContentTypes = []string{
	"application/javascript",
	"application/json",
	"application/wasm",
	"application/x-javascript",
	"application/x-tar",
	"image/svg+xml",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

func isCacheableOrCompressibleContentType(contentTypes []string, contentType string) bool {
	if contentType == "" {
		return false
	}

	for _, t := range contentTypes {
		if strings.Contains(contentType, t) {
			return true
		}
	}

	return false
}
