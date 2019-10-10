package http

import "net/http"

// Version returns a decorated version of the given handler that inject ETag
// and cache headers in order let browsers detect if they need to continue their
// request for a resource.
func Version(h http.Handler, etag string) http.Handler {
	if etag != "" {
		etag = etagHeaderValue(etag)
	}

	return &version{
		handler: h,
		etag:    etag,
	}
}

type version struct {
	handler http.Handler
	etag    string
}

func (v *version) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if v.etag == "" {
		v.handler.ServeHTTP(w, r)
		return
	}

	w.Header().Set("ETag", v.etag)
	w.Header().Set("Cache-Control", "no-cache")

	etag := r.Header.Get("If-None-Match")
	if etag == v.etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	v.handler.ServeHTTP(w, r)
}

func etagHeaderValue(etag string) string {
	return `"` + etag + `"`
}
