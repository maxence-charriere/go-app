package app

import "github.com/murlokswarm/app/key"

// HTMLConfig represents a configuration for an html document.
type HTMLConfig struct {
	// The title.
	Title string

	// The meta data.
	Metas []Meta

	// The css file paths to include.
	// Inludes all .css files in resources/css if nil no set.
	CSS []string

	// The javascript file paths to include.
	// Inludes all .js files in resources/js if not set.
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

// MouseEvent represents an onmouse event arg.
type MouseEvent struct {
	ClientX   float64
	ClientY   float64
	PageX     float64
	PageY     float64
	ScreenX   float64
	ScreenY   float64
	Button    int
	Detail    int
	AltKey    bool
	CtrlKey   bool
	MetaKey   bool
	ShiftKey  bool
	InnerText string
	Source    EventSource
}

// WheelEvent represents an onwheel event arg.
type WheelEvent struct {
	DeltaX    float64
	DeltaY    float64
	DeltaZ    float64
	DeltaMode DeltaMode
	Source    EventSource
}

// DeltaMode is an indication of the units of measurement for a delta value.
type DeltaMode uint64

// KeyboardEvent represents an onkey event arg.
type KeyboardEvent struct {
	CharCode  rune
	KeyCode   key.Code
	Location  key.Location
	AltKey    bool
	CtrlKey   bool
	MetaKey   bool
	ShiftKey  bool
	InnerText string
	Source    EventSource
}

// DragAndDropEvent represents an ondrop event arg.
type DragAndDropEvent struct {
	Files         []string
	Data          string
	DropEffect    string
	EffectAllowed string
	Source        EventSource
}

// EventSource represents a descriptor to an event source.
type EventSource struct {
	GoappID string
	CompoID string
	ID      string
	Class   string
	Data    map[string]string
	Value   string
}
