package http

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"sync"
)

// ManifestHandler is a handler that serves a manifest file for progressive
// webapp support.
type ManifestHandler struct {
	BackgroundColor string
	Name            string
	Orientation     string
	ShortName       string
	Scope           string
	StartURL        string
	ThemeColor      string

	once         sync.Once
	lastModified string
	manifest     []byte
}

func (h *ManifestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(h.init)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Last-Modified", lastModified)
	w.Write(h.manifest)
}

func (h *ManifestHandler) init() {
	buffer := bytes.Buffer{}
	writer := gzip.NewWriter(&buffer)

	enc := json.NewEncoder(writer)
	enc.SetIndent("", "    ")

	if err := enc.Encode(manifest{
		BackgroundColor: h.BackgroundColor,
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
		Name:                      h.Name,
		Orientation:               h.Orientation,
		PreferRelatedApplications: true,
		RelatedApplications:       []interface{}{},
		ShortName:                 h.ShortName,
		Scope:                     h.Scope,
		StartURL:                  h.StartURL,
		ThemeColor:                h.ThemeColor,
	}); err != nil {
		panic(err)
	}

	writer.Close()
	h.manifest = buffer.Bytes()
}

type manifest struct {
	BackgroundColor           string         `json:"background_color,omitempty"`
	Display                   string         `json:"display"`
	Icons                     []manifestIcon `json:"icons"`
	Name                      string         `json:"name"`
	Orientation               string         `json:"orientation"`
	PreferRelatedApplications bool           `json:"prefer_related_applications"`
	RelatedApplications       []interface{}  `json:"related_applications"`
	ShortName                 string         `json:"short_name"`
	Scope                     string         `json:"scope"`
	StartURL                  string         `json:"start_url"`
	ThemeColor                string         `json:"theme_color,omitempty"`
}

type manifestIcon struct {
	Sizes string `json:"sizes"`
	Src   string `json:"src"`
	Type  string `json:"type"`
}
