package html

import (
	"bytes"
	"html/template"
)

//go:generate go run gen.go
//go:generate go fmt

// PageConfig is the struct that describes a page.
type PageConfig struct {
	// The title.
	Title string

	// The meta data.
	Metas []Meta

	// The default component rendering.
	DefaultComponent template.HTML

	// Enables application default style.
	AppStyle bool

	// The CSS filenames to include.
	CSS []string

	// The app.js code that is included in the page..
	AppJS template.JS

	// The javascript filenames to include.
	Javascripts []string
}

// Meta represents a page metadata.
type Meta struct {
	Name      MetaName
	Content   string
	HTTPEquiv MetaHTTPEquiv
}

// MetaName represents a metadata name value.
type MetaName string

// Constants that define metadata name values.
const (
	ApplicationNameMeta MetaName = "application-name"
	AuthorMeta          MetaName = "author"
	DescriptionMeta     MetaName = "description"
	GeneratorMeta       MetaName = "generator"
	KeywordsMeta        MetaName = "keywords"
	ViewportMeta        MetaName = "viewport"
)

// MetaHTTPEquiv represents a metadata http-equiv value.
type MetaHTTPEquiv string

// Constants that define metadata http-equiv values.
const (
	ContentTypeMeta  MetaHTTPEquiv = "content-type"
	DefaultStyleMeta MetaHTTPEquiv = "default-style"
	RefreshMeta      MetaHTTPEquiv = "refresh"
)

// Page generate an HTML page from the given configuration.
func Page(config PageConfig) string {
	var buffer bytes.Buffer

	tmpl := template.Must(template.New("").Parse(pageTemplate))
	tmpl.Execute(&buffer, config)
	return buffer.String()
}
