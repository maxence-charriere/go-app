package http

import (
	"net/http"
	"time"

	"github.com/maxence-charriere/app/pkg/log"
)

func Watch(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Info("request served").
			T("path", r.URL.Path).
			T("content-type", w.Header().Get("Content-Type")).
			T("content-encoding", w.Header().Get("Content-Encoding")).
			T("content-lenght", w.Header().Get("Content-Length")).
			T("duration", time.Now().Sub(start).String())
	})
}
