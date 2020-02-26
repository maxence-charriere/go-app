package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/maxence-charriere/go-app/pkg/log"
)

// Manifest is a handler that serves a manifest file for progressive webapp
// support.
type Manifest struct {
	BackgroundColor string
	Name            string
	Orientation     string
	ShortName       string
	Scope           string
	StartURL        string
	ThemeColor      string

	once sync.Once
	body []byte
}

// CanHandle returns whether it can handle the given request.
func (m *Manifest) CanHandle(r *http.Request) bool {
	return r.URL.Path == "/manifest.json"
}

func (m *Manifest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.once.Do(m.init)

	w.Header().Set("Content-Length", strconv.Itoa(len(m.body)))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Last-Modified", lastModified)
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(m.body); err != nil {
		log.Error("writing response failed").
			T("error", err).
			T("path", r.URL.Path)
	}
}

func (m *Manifest) init() {
	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	if err := enc.Encode(manifest{
		BackgroundColor: m.BackgroundColor,
		Display:         "standalone",
		Icons: []manifestIcon{
			{
				Sizes: "192x192",
				Src:   "/icon-192.png",
				Type:  "image/png",
			},
			{
				Sizes: "512x512",
				Src:   "/icon-512.png",
				Type:  "image/png",
			},
		},
		Name:                      m.Name,
		Orientation:               m.Orientation,
		PreferRelatedApplications: true,
		RelatedApplications:       []interface{}{},
		ShortName:                 m.ShortName,
		Scope:                     m.Scope,
		StartURL:                  m.StartURL,
		ThemeColor:                m.ThemeColor,
	}); err != nil {
		log.Error("generating manifest.json failed").
			T("error", err).
			Panic()
	}

	m.body = b.Bytes()
}

type manifest struct {
	BackgroundColor           string         `json:"background_color"`
	Display                   string         `json:"display"`
	Icons                     []manifestIcon `json:"icons"`
	Name                      string         `json:"name"`
	Orientation               string         `json:"orientation"`
	PreferRelatedApplications bool           `json:"prefer_related_applications"`
	RelatedApplications       []interface{}  `json:"related_applications"`
	ShortName                 string         `json:"short_name"`
	Scope                     string         `json:"scope"`
	StartURL                  string         `json:"start_url"`
	ThemeColor                string         `json:"theme_color"`
}

type manifestIcon struct {
	Sizes string `json:"sizes"`
	Src   string `json:"src"`
	Type  string `json:"type"`
}
