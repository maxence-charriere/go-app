package app

// HTMLConfig represents a configuration for an html document.
type HTMLConfig struct {
	// The title.
	Title string

	// The meta data.
	Metas []Meta

	// The CSS filenames to include.
	// Inludes all files in resources/css if nil.
	CSS []string

	// The javascript filenames to include.
	// Inludes all files in resources/js if nil.
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
