//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type tag struct {
	Name          string
	Type          tagType
	Doc           string
	Attrs         []attr
	EventHandlers []eventHandler
}

type tagType int

const (
	parent tagType = iota
	privateParent
	selfClosing
)

var tags = []tag{
	{
		// A:
		Name: "A",
		Doc:  "defines a hyperlink.",
		Attrs: withGlobalAttrs(attrsByNames(
			"download",
			"href",
			"hreflang",
			"media",
			"ping",
			"rel",
			"target",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Abbr",
		Doc:           "defines an abbreviation or an acronym.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Address",
		Doc:           "defines contact information for the author/owner of a document.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Area",
		Type: selfClosing,
		Doc:  "defines an area inside an image-map.",
		Attrs: withGlobalAttrs(attrsByNames(
			"alt",
			"coords",
			"download",
			"href",
			"hreflang",
			"media",
			"rel",
			"shape",
			"target",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Article",
		Doc:           "defines an article.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Aside",
		Doc:           "defines content aside from the page content.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Audio",
		Doc:  "defines sound content.",
		Attrs: withGlobalAttrs(attrsByNames(
			"autoplay",
			"controls",
			"crossorigin",
			"loop",
			"muted",
			"preload",
			"src",
		)...),
		EventHandlers: withMediaEventHandlers(withGlobalEventHandlers()...),
	},

	// B:
	{
		Name:          "B",
		Doc:           "defines bold text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Base",
		Type: selfClosing,
		Doc:  "specifies the base URL/target for all relative URLs in a document.",
		Attrs: withGlobalAttrs(attrsByNames(
			"href",
			"target",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Bdi",
		Doc:           "isolates a part of text that might be formatted in a different direction from other text outside it.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Bdo",
		Doc:           "overrides the current text direction.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Blockquote",
		Doc:  "defines a section that is quoted from another source.",
		Attrs: withGlobalAttrs(attrsByNames(
			"cite",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "Body",
		Type:  privateParent,
		Doc:   "defines the document's body.",
		Attrs: withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"onafterprint",
			"onbeforeprint",
			"onbeforeunload",
			"onerror",
			"onhashchange",
			"onload",
			"onmessage",
			"onoffline",
			"ononline",
			"onpagehide",
			"onpageshow",
			"onpopstate",
			"onresize",
			"onstorage",
			"onunload",
		)...),
	},
	{
		Name:          "Br",
		Type:          selfClosing,
		Doc:           "defines a single line break.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Button",
		Doc:  "defines a clickable button.",
		Attrs: withGlobalAttrs(attrsByNames(
			"autofocus",
			"disabled",
			"form",
			"formaction",
			"formenctype",
			"formmethod",
			"formnovalidate",
			"formtarget",
			"name",
			"type",
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// C:
	{
		Name: "Canvas",
		Doc:  "is used to draw graphics on the fly.",
		Attrs: withGlobalAttrs(attrsByNames(
			"height",
			"width",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Caption",
		Doc:           "defines a table caption.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Cite",
		Doc:           "defines the title of a work.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Code",
		Doc:           "defines a piece of computer code.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Col",
		Type: selfClosing,
		Doc:  "specifies column properties for each column within a colgroup element.",
		Attrs: withGlobalAttrs(attrsByNames(
			"span",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "ColGroup",
		Doc:  "specifies a group of one or more columns in a table for formatting.",
		Attrs: withGlobalAttrs(attrsByNames(
			"span",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// D:
	{
		Name: "Data",
		Doc:  "links the given content with a machine-readable translation.",
		Attrs: withGlobalAttrs(attrsByNames(
			"value",
		)...),
	},
	{
		Name:          "DataList",
		Doc:           "specifies a list of pre-defined options for input controls.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Dd",
		Doc:           "defines a description/value of a term in a description list.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Del",
		Doc:  "defines text that has been deleted from a document.",
		Attrs: withGlobalAttrs(attrsByNames(
			"cite",
			"datetime",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Details",
		Doc:  "defines additional details that the user can view or hide.",
		Attrs: withGlobalAttrs(attrsByNames(
			"open",
		)...),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"ontoggle",
		)...),
	},
	{
		Name:          "Dfn",
		Doc:           "represents the defining instance of a term.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Dialog",
		Doc:  "defines a dialog box or window.",
		Attrs: withGlobalAttrs(attrsByNames(
			"open",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Div",
		Doc:           "defines a section in a document.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Dl",
		Doc:           "defines a description list.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Dt",
		Doc:           "defines a term/name in a description list.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// E:
	{
		Name: "Elem",
		Doc:  "represents an customizable HTML element.",
		Attrs: withGlobalAttrs(attrsByNames(
			"xmlns",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "ElemSelfClosing",
		Type: selfClosing,
		Doc:  "represents a self closing custom HTML element.",
		Attrs: withGlobalAttrs(attrsByNames(
			"xmlns",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Em",
		Doc:           "defines emphasized text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Embed",
		Type: selfClosing,
		Doc:  "defines a container for an external (non-HTML) application.",
		Attrs: withGlobalAttrs(attrsByNames(
			"height",
			"src",
			"type",
			"width",
		)...),
		EventHandlers: withMediaEventHandlers(withGlobalEventHandlers()...),
	},

	// F:
	{
		Name: "FieldSet",
		Doc:  "groups related elements in a form.",
		Attrs: withGlobalAttrs(attrsByNames(
			"disabled",
			"form",
			"name",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "FigCaption",
		Doc:           "defines a caption for a figure element.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Figure",
		Doc:           "specifies self-contained content.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Footer",
		Doc:           "defines a footer for a document or section.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Form",
		Doc:  "defines an HTML form for user input.",
		Attrs: withGlobalAttrs(attrsByNames(
			"accept-charset",
			"action",
			"autocomplete",
			"enctype",
			"method",
			"name",
			"novalidate",
			"target",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// H:
	{
		Name:          "H1",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "H2",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "H3",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "H4",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "H5",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "H6",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "Head",
		Doc:   "defines information about the document.",
		Attrs: withGlobalAttrs(attrsByNames()...),
	},
	{
		Name:          "Header",
		Doc:           "defines a header for a document or section.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Hr",
		Type:          selfClosing,
		Doc:           "defines a thematic change in the content.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "Html",
		Type:  privateParent,
		Doc:   "defines the root of an HTML document.",
		Attrs: withGlobalAttrs(),
	},

	// I:
	{
		Name:          "I",
		Doc:           "defines a part of text in an alternate voice or mood.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "IFrame",
		Doc:  "defines an inline frame.",
		Attrs: withGlobalAttrs(attrsByNames(
			"allow",
			"allowfullscreen",
			"allowpaymentrequest",
			"height",
			"name",
			"referrerpolicy",
			"sandbox",
			"src",
			"srcdoc",
			"width",
			"loading",
		)...),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"onload",
		)...,
		),
	},
	{
		Name: "Img",
		Type: selfClosing,
		Doc:  "defines an image.",
		Attrs: withGlobalAttrs(attrsByNames(
			"alt",
			"crossorigin",
			"fetchpriority",
			"height",
			"ismap",
			"sizes",
			"src",
			"srcset",
			"usemap",
			"width",
		)...),
		EventHandlers: withMediaEventHandlers(withGlobalEventHandlers(
			eventHandlersByName(
				"onload",
			)...,
		)...),
	},
	{
		Name: "Input",
		Type: selfClosing,
		Doc:  "defines an input control.",
		Attrs: withGlobalAttrs(attrsByNames(
			"accept",
			"alt",
			"autocomplete",
			"autofocus",
			"capture",
			"checked",
			"dirname",
			"disabled",
			"form",
			"formaction",
			"formenctype",
			"formmethod",
			"formnovalidate",
			"formtarget",
			"height",
			"list",
			"max",
			"maxlength",
			"min",
			"multiple",
			"name",
			"pattern",
			"placeholder",
			"readonly",
			"required",
			"size",
			"src",
			"step",
			"type",
			"value",
			"width",
		)...),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"onload",
		)...,
		),
	},
	{
		Name:          "Ins",
		Doc:           "defines a text that has been inserted into a document.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// K:
	{
		Name:          "Kbd",
		Doc:           "defines keyboard input.",
		Attrs:         withGlobalAttrs(attrsByNames()...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// L:
	{
		Name: "Label",
		Doc:  "defines a label for an input element.",
		Attrs: withGlobalAttrs(attrsByNames(
			"for",
			"form",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Legend",
		Doc:           "defines a caption for a fieldset element.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Li",
		Doc:  "defines a list item.",
		Attrs: withGlobalAttrs(attrsByNames(
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Link",
		Type: selfClosing,
		Doc:  "defines the relationship between a document and an external resource (most used to link to style sheets).",
		Attrs: withGlobalAttrs(attrsByNames(
			"as",
			"crossorigin",
			"fetchpriority",
			"href",
			"hreflang",
			"media",
			"rel",
			"sizes",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"onload",
		)...),
	},

	// M:
	{
		Name:          "Main",
		Doc:           "specifies the main content of a document.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Map",
		Doc:  "defines a client-side image-map.",
		Attrs: withGlobalAttrs(attrsByNames(
			"name",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Mark",
		Doc:           "defines marked/highlighted text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Meta",
		Type: selfClosing,
		Doc:  ".",
		Attrs: withGlobalAttrs(attrsByNames(
			"charset",
			"content",
			"http-equiv",
			"name",
			"property",
		)...),
	},
	{
		Name: "Meter",
		Doc:  "defines a scalar measurement within a known range (a gauge).",
		Attrs: withGlobalAttrs(attrsByNames(
			"form",
			"high",
			"low",
			"max",
			"min",
			"optimum",
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// N:
	{
		Name:          "Nav",
		Doc:           "defines navigation links.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "NoScript",
		Doc:   "defines an alternate content for users that do not support client-side scripts.",
		Attrs: withGlobalAttrs(attrsByNames()...),
	},

	// O:
	{
		Name: "Object",
		Doc:  "defines an embedded object.",
		Attrs: withGlobalAttrs(attrsByNames(
			"data",
			"form",
			"height",
			"name",
			"type",
			"usemap",
			"width",
		)...),
		EventHandlers: withMediaEventHandlers(withGlobalEventHandlers()...),
	},
	{
		Name: "Ol",
		Doc:  "defines an ordered list.",
		Attrs: withGlobalAttrs(attrsByNames(
			"reversed",
			"start",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "OptGroup",
		Doc:  "defines a group of related options in a drop-down list.",
		Attrs: withGlobalAttrs(attrsByNames(
			"disabled",
			"label",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Option",
		Doc:  "defines an option in a drop-down list.",
		Attrs: withGlobalAttrs(attrsByNames(
			"disabled",
			"label",
			"selected",
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Output",
		Doc:  ".",
		Attrs: withGlobalAttrs(attrsByNames(
			"for",
			"form",
			"name",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// P:
	{
		Name:          "P",
		Doc:           "defines a paragraph.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Param",
		Type: selfClosing,
		Doc:  "defines a parameter for an object.",
		Attrs: withGlobalAttrs(attrsByNames(
			"name",
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Picture",
		Doc:           "defines a container for multiple image resources.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Pre",
		Doc:           "defines preformatted text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Progress",
		Doc:  "represents the progress of a task.",
		Attrs: withGlobalAttrs(attrsByNames(
			"max",
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// Q:
	{
		Name: "Q",
		Doc:  "defines a short quotation.",
		Attrs: withGlobalAttrs(attrsByNames(
			"cite",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// R:
	{
		Name:          "Rp",
		Doc:           "defines what to show in browsers that do not support ruby annotations.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Rt",
		Doc:           "defines an explanation/pronunciation of characters (for East Asian typography).",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Ruby",
		Doc:           "defines a ruby annotation (for East Asian typography).",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// S:
	{
		Name:          "S",
		Doc:           "Defines text that is no longer correct.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Samp",
		Doc:           "defines sample output from a computer program.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Script",
		Doc:  "defines a client-side script.",
		Attrs: withGlobalAttrs(attrsByNames(
			"async",
			"charset",
			"crossorigin",
			"defer",
			"src",
			"type",
		)...),
		EventHandlers: eventHandlersByName("onload"),
	},
	{
		Name:          "Section",
		Doc:           "defines a section in a document.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Select",
		Doc:  "defines a drop-down list.",
		Attrs: withGlobalAttrs(attrsByNames(
			"autofocus",
			"disabled",
			"form",
			"multiple",
			"name",
			"required",
			"size",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Small",
		Doc:           "defines smaller text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Source",
		Type: selfClosing,
		Doc:  ".",
		Attrs: withGlobalAttrs(attrsByNames(
			"src",
			"srcset",
			"media",
			"sizes",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Span",
		Doc:           "defines a section in a document.",
		Attrs:         withGlobalAttrs(attrsByNames()...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Strong",
		Doc:           "defines important text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Style",
		Doc:  "defines style information for a document.",
		Attrs: withGlobalAttrs(attrsByNames(
			"media",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"onload",
		)...),
	},
	{
		Name:          "Sub",
		Doc:           "defines subscripted text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Summary",
		Doc:           "defines a visible heading for a details element.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Sup",
		Doc:           "defines superscripted text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// T:
	{
		Name:          "Table",
		Doc:           "defines a table.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "TBody",
		Doc:           "groups the body content in a table.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Td",
		Doc:  "defines a cell in a table.",
		Attrs: withGlobalAttrs(attrsByNames(
			"colspan",
			"headers",
			"rowspan",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "Template",
		Doc:   "defines a template.",
		Attrs: withGlobalAttrs(attrsByNames()...),
	},
	{
		Name: "Textarea",
		Doc:  "defines a multiline input control (text area).",
		Attrs: withGlobalAttrs(attrsByNames(
			"autofocus",
			"cols",
			"dirname",
			"disabled",
			"form",
			"maxlength",
			"name",
			"placeholder",
			"readonly",
			"required",
			"rows",
			"wrap",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "TFoot",
		Doc:           "groups the footer content in a table.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Th",
		Doc:  "defines a header cell in a table.",
		Attrs: withGlobalAttrs(attrsByNames(
			"abbr",
			"colspan",
			"headers",
			"rowspan",
			"scope",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "THead",
		Doc:           "groups the header content in a table",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Time",
		Doc:  "defines a date/time.",
		Attrs: withGlobalAttrs(attrsByNames(
			"datetime",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "Title",
		Doc:   "defines a title for the document.",
		Attrs: withGlobalAttrs(attrsByNames()...),
	},
	{
		Name:          "Tr",
		Doc:           "defines a row in a table.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// U:
	{
		Name:          "U",
		Doc:           "defines text that should be stylistically different from normal text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Ul",
		Doc:           "defines an unordered list.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// V:
	{
		Name:          "Var",
		Doc:           "defines a variable.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Video",
		Doc:  "defines a video or movie.",
		Attrs: withGlobalAttrs(attrsByNames(
			"autoplay",
			"controls",
			"crossorigin",
			"height",
			"loop",
			"muted",
			"poster",
			"preload",
			"src",
			"width",
		)...),
		EventHandlers: withMediaEventHandlers(withGlobalEventHandlers()...),
	},
	{
		Name:          "Wbr",
		Doc:           "defines a possible line-break.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
}

type attr struct {
	Name         string
	NameOverride string
	Type         string
	Key          bool
	Doc          string
}

var attrs = map[string]attr{
	// A:
	"abbr": {
		Name: "Abbr",
		Type: "fmt",
		Doc:  "specifies an abbreviated version of the content in a header cell with the given format and values.",
	},
	"accept": {
		Name: "Accept",
		Type: "fmt",
		Doc:  "specifies the types of files that the server accepts (only for file type) with the given format and values.",
	},
	"allow": {
		Name: "Allow",
		Type: "fmt",
		Doc:  "specifies a feature policy with the given format and values. Can be called multiple times to set multiple policies.",
	},
	"allowfullscreen": {
		Name: "AllowFullscreen",
		Type: "bool|force",
		Doc:  "reports whether an iframe can activate fullscreen mode.",
	},
	"allowpaymentrequest": {
		Name: "AllowPaymentRequest",
		Type: "bool|force",
		Doc:  "reports whether an iframe should be allowed to invoke the Payment Request API",
	},
	"aria-*": {
		Name: "Aria",
		Type: "aria|value",
		Doc:  "stores accessible rich internet applications (ARIA) data.",
	},
	"attribute": {
		Name: "Attr",
		Type: "attr|value",
		Doc:  "sets the named attribute with the given value.",
	},
	"accept-charset": {
		Name:         "AcceptCharset",
		NameOverride: "accept-charset",
		Type:         "string",
		Doc:          "specifies the character encodings that are to be used for the form submission.",
	},
	"accesskey": {
		Name: "AccessKey",
		Type: "fmt",
		Doc:  "specifies a shortcut key with the given format and values to activate/focus an element.",
	},
	"action": {
		Name: "Action",
		Type: "fmt",
		Doc:  "specifies where to send the form-data with the given format and values when a form is submitted.",
	},
	"alt": {
		Name: "Alt",
		Type: "fmt",
		Doc:  "specifies an alternate text with the given format and values when the original element fails to display.",
	},
	"as": {
		Name: "As",
		Type: "fmt",
		Doc:  "specifies a resource type to preload with the given format and values.",
	},
	"async": {
		Name: "Async",
		Type: "bool",
		Doc:  "specifies that the script is executed asynchronously (only for external scripts).",
	},
	"autocomplete": {
		Name: "AutoComplete",
		Type: "on/off",
		Doc:  "specifies whether the element should have autocomplete enabled.",
	},
	"autofocus": {
		Name: "AutoFocus",
		Type: "bool",
		Doc:  "specifies that the element should automatically get focus when the page loads.",
	},
	"autoplay": {
		Name: "AutoPlay",
		Type: "bool",
		Doc:  "specifies that the audio/video will start playing as soon as it is ready.",
	},

	// C:
	"capture": {
		Name: "Capture",
		Type: "fmt",
		Doc:  "specifies the capture input method in file upload controls with the given format and values.",
	},
	"charset": {
		Name: "Charset",
		Type: "fmt",
		Doc:  "specifies the character encoding with the given format and values.",
	},
	"checked": {
		Name: "Checked",
		Type: "bool",
		Doc:  "specifies that an input element should be pre-selected when the page loads (for checkbox or radio types).",
	},
	"cite": {
		Name: "Cite",
		Type: "fmt",
		Doc:  "specifies a URL which explains the quote/deleted/inserted text with the given format and values.",
	},
	"class": {
		Name: "Class",
		Type: "string|class",
		Doc:  "specifies one or more classnames for an element (refers to a class in a style sheet).",
	},
	"cols": {
		Name: "Cols",
		Type: "int",
		Doc:  "specifies the visible width of a text area.",
	},
	"colspan": {
		Name: "ColSpan",
		Type: "int",
		Doc:  "specifies the number of columns a table cell should span.",
	},
	"content": {
		Name: "Content",
		Type: "fmt",
		Doc:  "specifies the value associated with the http-equiv or name attribute using the given format and values.",
	},
	"contenteditable": {
		Name: "ContentEditable",
		Type: "bool",
		Doc:  "specifies whether the content of an element is editable or not.",
	},
	"controls": {
		Name: "Controls",
		Type: "bool",
		Doc:  "specifies that audio/video controls should be displayed (such as a play/pause button etc).",
	},
	"coords": {
		Name: "Coords",
		Type: "fmt",
		Doc:  "specifies the coordinates of the area with the given format and values.",
	},
	"crossorigin": {
		Name: "CrossOrigin",
		Type: "fmt",
		Doc:  "sets the mode of the request to an HTTP CORS Request with the given format and values.",
	},

	// D:
	"data": {
		Name: "Data",
		Type: "fmt",
		Doc:  "specifies the URL of the resource to be used by the object with the given format and values.",
	},
	"data-*": {
		Name: "DataSet",
		Type: "data|value",
		Doc:  "stores custom data private to the page or application.",
	},
	"datasets": {
		Name: "DataSets",
		Type: "data|map",
		Doc:  "specifies datsets for an element. Can be called multiple times to set multiple data set.",
	},
	"datetime": {
		Name: "DateTime",
		Type: "fmt",
		Doc:  "specifies the date and time with the given format and values.",
	},
	"default": {
		Name: "Default",
		Type: "bool",
		Doc:  "specifies that the track is to be enabled if the user's preferences do not indicate that another track would be more appropriate.",
	},
	"defer": {
		Name: "Defer",
		Type: "bool",
		Doc:  "specifies that the script is executed when the page has finished parsing (only for external scripts).",
	},
	"dir": {
		Name: "Dir",
		Type: "fmt",
		Doc:  "specifies the text direction for the content in an element with the given format and values.",
	},
	"dirname": {
		Name: "DirName",
		Type: "fmt",
		Doc:  "specifies that the text direction will be submitted using the given format and values.",
	},
	"disabled": {
		Name: "Disabled",
		Type: "bool",
		Doc:  "specifies that the specified element/group of elements should be disabled.",
	},
	"download": {
		Name: "Download",
		Type: "fmt",
		Doc:  "specifies that the target will be downloaded when a user clicks on the hyperlink. Uses the given format and values.",
	},
	"draggable": {
		Name: "Draggable",
		Type: "bool",
		Doc:  "specifies whether an element is draggable or not.",
	},

	// E:
	"enctype": {
		Name: "EncType",
		Type: "fmt",
		Doc:  "specifies how the form-data should be encoded when submitting it to the server (only for post method). Uses the given format and values.",
	},

	// F:
	"fetchpriority": {
		Name: "FetchPriority",
		Type: "string",
		Doc:  "specifies a hint given to the browser on how it should prioritize the fetch of the image relative to other images.",
	},
	"for": {
		Name: "For",
		Type: "fmt",
		Doc:  "specifies which form element(s) a label/calculation is bound to. Uses the given format and values.",
	},
	"form": {
		Name: "Form",
		Type: "fmt",
		Doc:  "specifies the name of the form the element belongs to. Uses the given format and values.",
	},
	"formaction": {
		Name: "FormAction",
		Type: "fmt",
		Doc:  "specifies where to send the form-data when a form is submitted. Only for submit type. Uses the given format and values.",
	},
	"formenctype": {
		Name: "FormEncType",
		Type: "fmt",
		Doc:  "specifies how form-data should be encoded before sending it to a server. Only for submit type. Uses the given format and values.",
	},
	"formmethod": {
		Name: "FormMethod",
		Type: "fmt",
		Doc:  "specifies how to send the form-data (which HTTP method to use). Only for submit type. Uses the given format and values.",
	},
	"formnovalidate": {
		Name: "FormNoValidate",
		Type: "bool",
		Doc:  "specifies that the form-data should not be validated on submission. Only for submit type.",
	},
	"formtarget": {
		Name: "FormTarget",
		Type: "fmt",
		Doc:  "specifies where to display the response after submitting the form. Only for submit type. Uses the given format and values.",
	},

	// H:
	"headers": {
		Name: "Headers",
		Type: "fmt",
		Doc:  "specifies one or more headers cells a cell is related to. Uses the given format and values.",
	},
	"height": {
		Name: "Height",
		Type: "int",
		Doc:  "specifies the height of the element (in pixels).",
	},
	"hidden": {
		Name: "Hidden",
		Type: "bool",
		Doc:  "specifies that an element is not yet, or is no longer relevant.",
	},
	"high": {
		Name: "High",
		Type: "float64",
		Doc:  "specifies the range that is considered to be a high value.",
	},
	"href": {
		Name: "Href",
		Type: "fmt",
		Doc:  "specifies the URL of the page the link goes to with the given format and values.",
	},
	"hreflang": {
		Name: "HrefLang",
		Type: "fmt",
		Doc:  "specifies the language of the linked document with the given format and values.",
	},
	"http-equiv": {
		Name:         "HTTPEquiv",
		NameOverride: "http-equiv",
		Type:         "string",
		Doc:          "provides an HTTP header for the information/value of the content attribute.",
	},

	// I:
	"id": {
		Name: "ID",
		Type: "fmt",
		Doc:  "specifies a unique id for an element with the given format and values.",
	},
	"ismap": {
		Name: "IsMap",
		Type: "bool",
		Doc:  "specifies an image as a server-side image-map.",
	},

	// K:
	"kind": {
		Name: "Kind",
		Type: "fmt",
		Doc:  "specifies the kind of text track with the given format and values.",
	},

	// L:
	"label": {
		Name: "Label",
		Type: "fmt",
		Doc:  "specifies a shorter label for the option with the given format and values.",
	},
	"lang": {
		Name: "Lang",
		Type: "fmt",
		Doc:  "specifies the language of the element's content with the given format and values.",
	},
	"list": {
		Name: "List",
		Type: "fmt",
		Doc:  "refers to a datalist element that contains pre-defined options for an input element. Uses the given format and values.",
	},
	"loading": {
		Name: "Loading",
		Type: "fmt",
		Doc:  "indicates how the browser should load the iframe (eager|lazy). Uses the given format and values.",
	},
	"loop": {
		Name: "Loop",
		Type: "bool",
		Doc:  "specifies that the audio/video will start over again, every time it is finished.",
	},
	"low": {
		Name: "Low",
		Type: "float64",
		Doc:  "specifies the range that is considered to be a low value.",
	},

	// M:
	"max": {
		Name: "Max",
		Type: "any",
		Doc:  "Specifies the maximum value.",
	},
	"maxlength": {
		Name: "MaxLength",
		Type: "int",
		Doc:  "specifies the maximum number of characters allowed in an element.",
	},
	"media": {
		Name: "Media",
		Type: "fmt",
		Doc:  "specifies what media/device the linked document is optimized for. Uses the given format and values.",
	},
	"method": {
		Name: "Method",
		Type: "fmt",
		Doc:  "specifies the HTTP method to use when sending form-data. Uses the given format and values.",
	},
	"min": {
		Name: "Min",
		Type: "any",
		Doc:  "specifies a minimum value.",
	},
	"multiple": {
		Name: "Multiple",
		Type: "bool",
		Doc:  "specifies that a user can enter more than one value.",
	},
	"muted": {
		Name: "Muted",
		Type: "bool",
		Doc:  "specifies that the audio output of the video should be muted.",
	},

	// N:
	"name": {
		Name: "Name",
		Type: "fmt",
		Doc:  "specifies the name of the element with the given format and values.",
	},
	"novalidate": {
		Name: "NoValidate",
		Type: "bool",
		Doc:  "specifies that the form should not be validated when submitted.",
	},

	// O:
	"open": {
		Name: "Open",
		Type: "bool",
		Doc:  "specifies that the details should be visible (open) to the user.",
	},
	"optimum": {
		Name: "Optimum",
		Type: "float64",
		Doc:  "specifies what value is the optimal value for the gauge.",
	},

	// P:
	"pattern": {
		Name: "Pattern",
		Type: "fmt",
		Doc:  "specifies a regular expression that an input element's value is checked against. Uses the given format and values.",
	},
	"ping": {
		Name: "Ping",
		Type: "fmt",
		Doc:  "specifies a list of URLs to be notified if the user follows the hyperlink. Uses the given format and values.",
	},
	"placeholder": {
		Name: "Placeholder",
		Type: "fmt",
		Doc:  "specifies a short hint that describes the expected value of the element. Uses the given format and values.",
	},
	"poster": {
		Name: "Poster",
		Type: "fmt",
		Doc:  "specifies an image to be shown while the video is downloading, or until the user hits the play button. Uses the given format and values.",
	},
	"preload": {
		Name: "Preload",
		Type: "fmt",
		Doc:  "specifies if and how the author thinks the audio/video should be loaded when the page loads. Uses the given format and values.",
	},
	"property": {
		Name: "Property",
		Type: "fmt",
		Doc:  "specifies the property name with the given format and values.",
	},

	// R:
	"readonly": {
		Name: "ReadOnly",
		Type: "bool",
		Doc:  "specifies that the element is read-only.",
	},
	"referrerpolicy": {
		Name: "ReferrerPolicy",
		Type: "fmt",
		Doc:  "specifies how much/which referrer information that will be sent when processing the iframe attributes. Uses the given format and values.",
	},
	"rel": {
		Name: "Rel",
		Type: "fmt",
		Doc:  "specifies the relationship between the current document and the linked document. uses the given format and values.",
	},
	"required": {
		Name: "Required",
		Type: "bool",
		Doc:  "specifies that the element must be filled out before submitting the form.",
	},
	"reversed": {
		Name: "Reversed",
		Type: "bool",
		Doc:  "specifies that the list order should be descending (9,8,7...).",
	},
	"role": {
		Name: "Role",
		Type: "fmt",
		Doc:  "specifies to parsing software the exact function of an element (and its children). Uses the given format and values.",
	},
	"rows": {
		Name: "Rows",
		Type: "int",
		Doc:  "specifies the visible number of lines in a text area.",
	},
	"rowspan": {
		Name: "Rowspan",
		Type: "int",
		Doc:  "specifies the number of rows a table cell should span.",
	},

	// S:
	"sandbox": {
		Name: "Sandbox",
		Type: "any",
		Doc:  "enables an extra set of restrictions for the content in an iframe.",
	},
	"scope": {
		Name: "Scope",
		Type: "fmt",
		Doc:  "specifies whether a header cell is a header for a column, row, or group of columns or rows. Uses the given format and values.",
	},
	"selected": {
		Name: "Selected",
		Type: "bool",
		Doc:  "specifies that an option should be pre-selected when the page loads.",
	},
	"shape": {
		Name: "Shape",
		Type: "fmt",
		Doc:  "specifies the shape of the area with the given format and values.",
	},
	"size": {
		Name: "Size",
		Type: "int",
		Doc:  "specifies the width.",
	},
	"sizes": {
		Name: "Sizes",
		Type: "fmt",
		Doc:  "specifies the size of the linked resource with the given format and values.",
	},
	"span": {
		Name: "Span",
		Type: "int",
		Doc:  "specifies the number of columns to span.",
	},
	"spellcheck": {
		Name: "Spellcheck",
		Type: "bool|force",
		Doc:  "specifies whether the element is to have its spelling and grammar checked or not.",
	},
	"src": {
		Name: "Src",
		Type: "fmt",
		Doc:  "specifies the URL of the media file with the given format and values.",
	},
	"srcdoc": {
		Name: "SrcDoc",
		Type: "fmt",
		Doc:  "specifies the HTML content of the page to show in the iframe with the given format and values.",
	},
	"srclang": {
		Name: "SrcLang",
		Type: "fmt",
		Doc:  `specifies the language of the track text data (required if kind = "subtitles"). Uses the given format and values.`,
	},
	"srcset": {
		Name: "SrcSet",
		Type: "fmt",
		Doc:  "specifies the URL of the image to use in different situations with the given format and values.",
	},
	"start": {
		Name: "Start",
		Type: "int",
		Doc:  "specifies the start value of the ordered list.",
	},
	"step": {
		Name: "Step",
		Type: "float64",
		Doc:  "specifies the legal number intervals for an input field.",
	},
	"style": {
		Name: "Style",
		Type: "style",
		Doc:  "specifies a CSS style for an element. Can be called multiple times to set multiple css styles.",
	},
	"styles": {
		Name: "Styles",
		Type: "style|map",
		Doc:  "specifies CSS styles for an element. Can be called multiple times to set multiple css styles.",
	},

	// T:
	"tabindex": {
		Name: "TabIndex",
		Type: "int",
		Doc:  "specifies the tabbing order of an element.",
	},
	"target": {
		Name: "Target",
		Type: "fmt",
		Doc:  "specifies the target for where to open the linked document or where to submit the form. Uses the given format and values.",
	},
	"title": {
		Name: "Title",
		Type: "fmt",
		Doc:  "specifies extra information about an element with the given format and values.",
	},
	"type": {
		Name: "Type",
		Type: "fmt",
		Doc:  "specifies the type of element with the given format and values.",
	},

	// U:
	"usemap": {
		Name: "UseMap",
		Type: "fmt",
		Doc:  "specifies an image as a client-side image-map. Uses the given format and values.",
	},

	// V:
	"value": {
		Name: "Value",
		Type: "any",
		Doc:  "specifies the value of the element.",
	},

	// W:
	"width": {
		Name: "Width",
		Type: "int",
		Doc:  "specifies the width of the element.",
	},
	"wrap": {
		Name: "Wrap",
		Type: "fmt",
		Doc:  "specifies how the text in a text area is to be wrapped when submitted in a form. Uses the given format and values.",
	},
	"xmlns": {
		Name: "XMLNS",
		Type: "xmlns",
		Doc:  "specifies the xml namespace of the element.",
	},
}

func attrsByNames(names ...string) []attr {
	res := make([]attr, 0, len(names))
	for _, n := range names {
		attr, ok := attrs[n]
		if !ok {
			panic("unkowmn attr: " + n)
		}
		res = append(res, attr)
	}

	sort.Slice(res, func(i, j int) bool {
		return strings.Compare(res[i].Name, res[j].Name) <= 0
	})

	return res
}

func withGlobalAttrs(attrs ...attr) []attr {
	attrs = append(attrs, attrsByNames(
		"accesskey",
		"aria-*",
		"class",
		"contenteditable",
		"data-*",
		"datasets",
		"dir",
		"draggable",
		"hidden",
		"id",
		"lang",
		"role",
		"spellcheck",
		"style",
		"styles",
		"tabindex",
		"title",
		"attribute",
	)...)

	sort.Slice(attrs, func(i, j int) bool {
		return strings.Compare(attrs[i].Name, attrs[j].Name) <= 0
	})

	return attrs
}

type eventHandler struct {
	Name string
	Doc  string
}

var eventHandlers = map[string]eventHandler{
	// Window events:
	"onafterprint": {
		Name: "OnAfterPrint",
		Doc:  "runs the given handler after the document is printed.",
	},
	"onbeforeprint": {
		Name: "OnBeforePrint",
		Doc:  "calls the given handler before the document is printed.",
	},
	"onbeforeunload": {
		Name: "OnBeforeUnload",
		Doc:  "calls the given handler when the document is about to be unloaded.",
	},
	"onerror": {
		Name: "OnError",
		Doc:  "calls the given handler when an error occurs.",
	},
	"onhashchange": {
		Name: "OnHashChange",
		Doc:  "calls the given handler when there has been changes to the anchor part of the a URL.",
	},
	"onload": {
		Name: "OnLoad",
		Doc:  "calls the given handler after the element is finished loading.",
	},
	"onmessage": {
		Name: "OnMessage",
		Doc:  "calls then given handler when a message is triggered.",
	},
	"onoffline": {
		Name: "OnOffline",
		Doc:  "calls the given handler when the browser starts to work offline.",
	},
	"ononline": {
		Name: "OnOnline",
		Doc:  "calls the given handler when the browser starts to work online.",
	},
	"onpagehide": {
		Name: "OnPageHide",
		Doc:  "calls the given handler when a user navigates away from a page.",
	},
	"onpageshow": {
		Name: "OnPageShow",
		Doc:  "calls the given handler when a user navigates to a page.",
	},
	"onpopstate": {
		Name: "OnPopState",
		Doc:  "calls the given handler when the window's history changes.",
	},
	"onresize": {
		Name: "OnResize",
		Doc:  "calls the given handler when the browser window is resized.",
	},
	"onstorage": {
		Name: "OnStorage",
		Doc:  "calls the given handler when a Web Storage area is updated.",
	},
	"onunload": {
		Name: "OnUnload",
		Doc:  "calls the given handler once a page has unloaded (or the browser window has been closed).",
	},

	// Form events:
	"onblur": {
		Name: "OnBlur",
		Doc:  "calls the given handler when the element loses focus.",
	},
	"onchange": {
		Name: "OnChange",
		Doc:  "calls the given handler when the value of the element is changed.",
	},
	"oncontextmenu": {
		Name: "OnContextMenu",
		Doc:  "calls the given handler when a context menu is triggered.",
	},
	"onfocus": {
		Name: "OnFocus",
		Doc:  "calls the given handler when the element gets focus.",
	},
	"oninput": {
		Name: "OnInput",
		Doc:  "calls the given handler when an element gets user input.",
	},
	"oninvalid": {
		Name: "OnInvalid",
		Doc:  "calls the given handler when an element is invalid.",
	},
	"onreset": {
		Name: "OnReset",
		Doc:  "calls the given handler when the Reset button in a form is clicked.",
	},
	"onsearch": {
		Name: "OnSearch",
		Doc:  `calls the given handler when the user writes something in a search field.`,
	},
	"onselect": {
		Name: "OnSelect",
		Doc:  "calls the given handler after some text has been selected in an element.",
	},
	"onsubmit": {
		Name: "OnSubmit",
		Doc:  "calls the given handler when a form is submitted.",
	},

	// Keyboard events:
	"onkeydown": {
		Name: "OnKeyDown",
		Doc:  "calls the given handler when a user is pressing a key.",
	},
	"onkeypress": {
		Name: "OnKeyPress",
		Doc:  "calls the given handler when a user presses a key.",
	},
	"onkeyup": {
		Name: "OnKeyUp",
		Doc:  "calls the given handler when a user releases a key.",
	},

	// Mouse events:
	"onclick": {
		Name: "OnClick",
		Doc:  "calls the given handler when there is a mouse click on the element.",
	},
	"ondblclick": {
		Name: "OnDblClick",
		Doc:  "calls the given handler when there is a mouse double-click on the element.",
	},
	"onmousedown": {
		Name: "OnMouseDown",
		Doc:  "calls the given handler when a mouse button is pressed down on an element.",
	},
	"onmouseenter": {
		Name: "OnMouseEnter",
		Doc:  "calls the given handler when a mouse button is initially moved so that its hotspot is within the element at which the event was fired.",
	},
	"onmouseleave": {
		Name: "OnMouseLeave",
		Doc:  "calls the given handler when the mouse pointer is fired when the pointer has exited the element and all of its descendants.",
	},
	"onmousemove": {
		Name: "OnMouseMove",
		Doc:  "calls the given handler when the mouse pointer is moving while it is over an element.",
	},
	"onmouseout": {
		Name: "OnMouseOut",
		Doc:  "calls the given handler when the mouse pointer moves out of an element.",
	},
	"onmouseover": {
		Name: "OnMouseOver",
		Doc:  "calls the given handler when the mouse pointer moves over an element.",
	},
	"onmouseup": {
		Name: "OnMouseUp",
		Doc:  "calls the given handler when a mouse button is released over an element.",
	},
	"onwheel": {
		Name: "OnWheel",
		Doc:  "calls the given handler when the mouse wheel rolls up or down over an element.",
	},

	// Drag events:
	"ondrag": {
		Name: "OnDrag",
		Doc:  "calls the given handler when an element is dragged.",
	},
	"ondragend": {
		Name: "OnDragEnd",
		Doc:  "calls the given handler at the end of a drag operation.",
	},
	"ondragenter": {
		Name: "OnDragEnter",
		Doc:  "calls the given handler when an element has been dragged to a valid drop target.",
	},
	"ondragleave": {
		Name: "OnDragLeave",
		Doc:  "calls the given handler when an element leaves a valid drop target.",
	},
	"ondragover": {
		Name: "OnDragOver",
		Doc:  "calls the given handler when an element is being dragged over a valid drop target.",
	},
	"ondragstart": {
		Name: "OnDragStart",
		Doc:  "calls the given handler at the start of a drag operation.",
	},
	"ondrop": {
		Name: "OnDrop",
		Doc:  "calls the given handler when dragged element is being dropped.",
	},
	"onscroll": {
		Name: "OnScroll",
		Doc:  "calls the given handler when an element's scrollbar is being scrolled.",
	},

	// Clipboard event:
	"oncopy": {
		Name: "OnCopy",
		Doc:  "calls the given handler when the user copies the content of an element.",
	},
	"oncut": {
		Name: "OnCut",
		Doc:  "calls the given handler when the user cuts the content of an element.",
	},
	"onpaste": {
		Name: "OnPaste",
		Doc:  "calls the given handler when the user pastes some content in an element.",
	},

	// Media events:
	"onabort": {
		Name: "OnAbort",
		Doc:  "calls the given handler on abort.",
	},
	"oncanplay": {
		Name: "OnCanPlay",
		Doc:  "calls the given handler when a file is ready to start playing (when it has buffered enough to begin).",
	},
	"oncanplaythrough": {
		Name: "OnCanPlayThrough",
		Doc:  "calls the given handler when a file can be played all the way to the end without pausing for buffering.",
	},
	"oncuechange": {
		Name: "OnCueChange",
		Doc:  "calls the given handler when the cue changes in a track element.",
	},
	"ondurationchange": {
		Name: "OnDurationChange",
		Doc:  "calls the given handler when the length of the media changes.",
	},
	"onemptied": {
		Name: "OnEmptied",
		Doc:  "calls the given handler when something bad happens and the file is suddenly unavailable (like unexpectedly disconnects).",
	},
	"onended": {
		Name: "OnEnded",
		Doc:  "calls the given handler when the media has reach the end.",
	},
	"onloadeddata": {
		Name: "OnLoadedData",
		Doc:  "calls the given handler when media data is loaded.",
	},
	"onloadedmetadata": {
		Name: "OnLoadedMetaData",
		Doc:  "calls the given handler when meta data (like dimensions and duration) are loaded.",
	},
	"onloadstart": {
		Name: "OnLoadStart",
		Doc:  "calls the given handler just as the file begins to load before anything is actually loaded.",
	},
	"onpause": {
		Name: "OnPause",
		Doc:  "calls the given handler when the media is paused either by the user or programmatically.",
	},
	"onplay": {
		Name: "OnPlay",
		Doc:  "calls the given handler when the media is ready to start playing.",
	},
	"onplaying": {
		Name: "OnPlaying",
		Doc:  "calls the given handler when the media actually has started playing.",
	},
	"onprogress": {
		Name: "OnProgress",
		Doc:  "calls the given handler when the browser is in the process of getting the media data.",
	},
	"onratechange": {
		Name: "OnRateChange",
		Doc:  "calls the given handler each time the playback rate changes (like when a user switches to a slow motion or fast forward mode).",
	},
	"onseeked": {
		Name: "OnSeeked",
		Doc:  "calls the given handler when the seeking attribute is set to false indicating that seeking has ended.",
	},
	"onseeking": {
		Name: "OnSeeking",
		Doc:  "calls the given handler when the seeking attribute is set to true indicating that seeking is active.",
	},
	"onstalled": {
		Name: "OnStalled",
		Doc:  "calls the given handler when the browser is unable to fetch the media data for whatever reason.",
	},
	"onsuspend": {
		Name: "OnSuspend",
		Doc:  "calls the given handler when fetching the media data is stopped before it is completely loaded for whatever reason.",
	},
	"ontimeupdate": {
		Name: "OnTimeUpdate",
		Doc:  "calls the given handler when the playing position has changed (like when the user fast forwards to a different point in the media).",
	},
	"onvolumechange": {
		Name: "OnVolumeChange",
		Doc:  `calls the given handler each time the volume is changed which (includes setting the volume to "mute").`,
	},
	"onwaiting": {
		Name: "OnWaiting",
		Doc:  "calls the given handler when the media has paused but is expected to resume (like when the media pauses to buffer more data).",
	},

	// Miscs events:
	"ontoggle": {
		Name: "OnToggle",
		Doc:  "calls the given handler when the user opens or closes the details element.",
	},
}

func eventHandlersByName(names ...string) []eventHandler {
	res := make([]eventHandler, 0, len(names))
	for _, n := range names {
		h, ok := eventHandlers[n]
		if !ok {
			panic("unknown event handler: " + n)
		}
		res = append(res, h)
	}

	sort.Slice(res, func(i, j int) bool {
		return strings.Compare(res[i].Name, res[j].Name) <= 0
	})

	return res
}

func withGlobalEventHandlers(handlers ...eventHandler) []eventHandler {
	handlers = append(handlers, eventHandlersByName(
		"onblur",
		"onchange",
		"oncontextmenu",
		"onfocus",
		"oninput",
		"oninvalid",
		"onreset",
		"onsearch",
		"onselect",
		"onsubmit",

		"onkeydown",
		"onkeypress",
		"onkeyup",

		"onclick",
		"ondblclick",
		"onmousedown",
		"onmouseenter",
		"onmouseleave",
		"onmousemove",
		"onmouseout",
		"onmouseover",
		"onmouseup",
		"onwheel",

		"ondrag",
		"ondragend",
		"ondragenter",
		"ondragleave",
		"ondragover",
		"ondragstart",
		"ondrop",
		"onscroll",

		"oncopy",
		"oncut",
		"onpaste",
	)...)

	sort.Slice(handlers, func(i, j int) bool {
		return strings.Compare(handlers[i].Name, handlers[j].Name) <= 0
	})

	return handlers
}

func withMediaEventHandlers(handlers ...eventHandler) []eventHandler {
	handlers = append(handlers, eventHandlersByName(
		"onabort",
		"oncanplay",
		"oncanplaythrough",
		"oncuechange",
		"ondurationchange",
		"onemptied",
		"onended",
		"onerror",
		"onloadeddata",
		"onloadedmetadata",
		"onloadstart",
		"onpause",
		"onplay",
		"onplaying",
		"onprogress",
		"onratechange",
		"onseeked",
		"onseeking",
		"onstalled",
		"onsuspend",
		"ontimeupdate",
		"onvolumechange",
		"onwaiting",
	)...)

	sort.Slice(handlers, func(i, j int) bool {
		return strings.Compare(handlers[i].Name, handlers[j].Name) <= 0
	})

	return handlers
}

func main() {
	generateHTMLGo()
	generateHTMLTestGo()
}

func generateHTMLGo() {
	f, err := os.Create("html_gen.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "package app")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(f, `
import (
	"strings"
)
		`)

	for _, t := range tags {
		writeInterface(f, t)

		switch t.Name {
		case "Elem", "ElemSelfClosing":
			fmt.Fprintf(f, `
			// %s returns an HTML element that %s
			func %s(tag string) HTML%s {
				e := &html%s{
					htmlElement: htmlElement{
						tag: tag,
						isSelfClosing: %v,
					},
				}

				return e
			}
			`,
				t.Name,
				t.Doc,
				t.Name,
				t.Name,
				t.Name,
				t.Type == selfClosing,
			)

		default:
			fmt.Fprintf(f, `
			// %s returns an HTML element that %s
			func %s() HTML%s {
				e := &html%s{
					htmlElement: htmlElement{
						tag: "%s",
						isSelfClosing: %v,
					},
				}

				return e
			}
			`,
				t.Name,
				t.Doc,
				t.Name,
				t.Name,
				t.Name,
				strings.ToLower(t.Name),
				t.Type == selfClosing,
			)
		}

		fmt.Fprintln(f)
		fmt.Fprintln(f)
		writeStruct(f, t)
		fmt.Fprintln(f)
		fmt.Fprintln(f)
	}

}

func writeInterface(w io.Writer, t tag) {
	fmt.Fprintf(w, `
		// HTML%s is the interface that describes a "%s" HTML element.
		type HTML%s interface {
			UI
		`,
		t.Name,
		strings.ToLower(t.Name),
		t.Name,
	)

	switch t.Type {
	case parent:
		fmt.Fprintf(w, `
			// Body set the content of the element.
			Body(elems ...UI) HTML%s 
		`, t.Name)

		fmt.Fprintf(w, `
			// Text sets the content of the element with a text node containing the stringified given value.
			Text(v any) HTML%s
		`, t.Name)

		fmt.Fprintf(w, `
			// Textf sets the content of the element with the given format and values.
			Textf(format string, v ...any) HTML%s
		`, t.Name)

	case privateParent:
		fmt.Fprintf(w, `
			privateBody(elems ...UI) HTML%s 
		`, t.Name)
	}

	for _, a := range t.Attrs {
		fmt.Fprintln(w)
		fmt.Fprintln(w)

		fmt.Fprintf(w, "// %s %s\n", a.Name, a.Doc)
		writeAttrFunction(w, a, t, true)
	}

	fmt.Fprintln(w)

	fmt.Fprintf(w, `
		// On registers the given event handler to the specified event.
		On(event string, h EventHandler, scope ...any) HTML%s 
	`, t.Name)

	for _, e := range t.EventHandlers {
		fmt.Fprintln(w)
		fmt.Fprintln(w)

		fmt.Fprintf(w, "// %s %s\n", e.Name, e.Doc)
		writeEventFunction(w, e, t, true)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "}")
}

func writeStruct(w io.Writer, t tag) {
	fmt.Fprintf(w, `type html%s struct {
			htmlElement
		}`, t.Name)

	switch t.Type {
	case parent:
		fmt.Fprintf(w, `
			func (e *html%s) Body(v ...UI) HTML%s {
				e.setChildren(v...)
				return e
			}
			`,
			t.Name,
			t.Name,
		)

		if t.Name == "Textarea" {
			fmt.Fprintf(w, `
			func (e *html%s) Text(v any) HTML%s {
				e.setAttr("value", v)
				return e
			}
			`,
				t.Name,
				t.Name,
			)
			fmt.Fprintf(w, `
			func (e *html%s) Textf(format string, v ...any) HTML%s {
				e.setAttr("value", FormatString(format, v...))
				return e
			}
			`,
				t.Name,
				t.Name,
			)
		} else {
			fmt.Fprintf(w, `
			func (e *html%s) Text(v any) HTML%s {
				return e.Body(Text(v))
			}
			`,
				t.Name,
				t.Name,
			)

			fmt.Fprintf(w, `
			func (e *html%s) Textf(format string, v ...any) HTML%s {
				return e.Body(Textf(format, v...))
			}
			`,
				t.Name,
				t.Name,
			)
		}

	case privateParent:
		fmt.Fprintf(w, `
			func (e *html%s) privateBody(v ...UI) HTML%s {
				e.setChildren(v...)
				return e
			}
			`,
			t.Name,
			t.Name,
		)
	}

	for _, a := range t.Attrs {
		fmt.Fprintln(w)
		fmt.Fprintln(w)

		writeAttrFunction(w, a, t, false)
	}

	fmt.Fprintln(w)

	fmt.Fprintf(w, `
		func (e *html%s) On(event string, h EventHandler, scope ...any)  HTML%s {
			e.setEventHandler(event, h, scope...)
			return e
		}
		`,
		t.Name,
		t.Name,
	)

	for _, e := range t.EventHandlers {
		fmt.Fprintln(w)
		fmt.Fprintln(w)

		writeEventFunction(w, e, t, false)
	}
}

func writeAttrFunction(w io.Writer, a attr, t tag, isInterface bool) {
	if !isInterface {
		fmt.Fprintf(w, "func (e *html%s)", t.Name)
	}

	var attrName string
	if a.NameOverride != "" {
		attrName = strings.ToLower(a.NameOverride)
	} else {
		attrName = strings.ToLower(a.Name)
	}

	switch a.Type {
	case "data|value":
		fmt.Fprintf(w, `%s(k string, v any) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("data-"+k, FormatString("%s", v))
				return e
			}`, "%v")
		}

	case "data|map":
		fmt.Fprintf(w, `%s(ds map[string]any) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				for k, v := range ds {
					e.DataSet(k, v)
				}
				return e
			}`)
		}

	case "attr|value":
		fmt.Fprintf(w, `%s(n string, v any) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr(n, v)
				return e
			}`)
		}

	case "aria|value":
		fmt.Fprintf(w, `%s(k string, v any) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("aria-"+k, FormatString("%s", v))
				return e
			}`, "%v")
		}

	case "style":
		fmt.Fprintf(w, `%s(k, v string) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("style", k+":"+v)
				return e
			}`)
		}

	case "style|map":
		fmt.Fprintf(w, `%s(s map[string]string) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				for k, v := range s {
					e.Style(k, v)
				}
				return e
			}`)
		}

	case "on/off":
		fmt.Fprintf(w, `%s(v bool) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				s := "off"
				if (v) {
					s = "on"
				}
	
				e.setAttr("%s", s)
				return e
			}`, attrName)
		}

	case "bool|force":
		fmt.Fprintf(w, `%s(v bool) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				s := "false"
				if (v) {
					s = "true"
				}
	
				e.setAttr("%s", s)
				return e
			}`, attrName)
		}

	case "string|class":
		fmt.Fprintf(w, `%s(v ...string) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("%s", strings.Join(v, " "))
				return e
			}`, attrName)
		}

	case "xmlns":
		fmt.Fprintf(w, `%s(v string) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintln(w, `{
				e.xmlns = v
				return e
			}`)
		}

	case "fmt":
		fmt.Fprintf(w, `%s(format string, v ...any) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("%s", FormatString(format, v...))
				return e
			}`, attrName)
		}

	default:
		fmt.Fprintf(w, `%s(v %s) HTML%s`, a.Name, a.Type, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("%s", v)
				return e
			}`, attrName)
		}
	}
}

func writeEventFunction(w io.Writer, e eventHandler, t tag, isInterface bool) {
	if !isInterface {
		fmt.Fprintf(w, `func (e *html%s)`, t.Name)
	}

	fmt.Fprintf(w, `%s (h EventHandler, scope ...any) HTML%s`, e.Name, t.Name)
	if !isInterface {
		fmt.Fprintf(w, `{
			e.setEventHandler("%s", h, scope...)
			return e
		}`, strings.TrimPrefix(strings.ToLower(e.Name), "on"))
	}
}

func generateHTMLTestGo() {
	f, err := os.Create("html_gen_test.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "package app")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(f, `
import (
	"testing"
)
		`)

	for _, t := range tags {
		fmt.Fprintln(f)
		fmt.Fprintf(f, `func Test%s(t *testing.T) {`, t.Name)
		fmt.Fprintln(f)

		switch t.Name {
		case "Elem", "ElemSelfClosing":
			fmt.Fprintf(f, `elem := %s("div")`, t.Name)

		default:
			fmt.Fprintf(f, `elem := %s()`, t.Name)
		}

		fmt.Fprintln(f)

		for _, a := range t.Attrs {
			fmt.Fprintf(f, `elem.%s(`, a.Name)

			switch a.Type {
			case "data|value", "aria|value", "attr|value":
				fmt.Fprintln(f, `"foo", "bar")`)

			case "data|map":
				fmt.Fprintln(f, `map[string]any{"foo": "bar"})`)

			case "style":
				fmt.Fprintln(f, `"color", "deepskyblue")`)

			case "style|map":
				fmt.Fprintln(f, `map[string]string{"color": "pink"})`)

			case "bool", "bool|force", "on/off":
				fmt.Fprintln(f, `true)`)
				fmt.Fprintf(f, `elem.%s(false)`, a.Name)
				fmt.Fprintln(f)

			case "int":
				fmt.Fprintln(f, `42)`)

			case "string":
				fmt.Fprintln(f, `"foo")`)

			case "url":
				fmt.Fprintln(f, `"http://foo.com")`)

			case "string|class":
				fmt.Fprintln(f, `"foo bar")`)

			case "xmlns":
				fmt.Fprintln(f, `"http://www.w3.org/2000/svg")`)

			case "fmt":
				fmt.Fprintln(f, `"hello %v", 42)`)

			default:
				fmt.Fprintln(f, `42)`)
			}
		}

		if len(t.EventHandlers) != 0 {
			fmt.Fprint(f, `
				h := func(ctx Context, e Event) {}
			`)
			fmt.Fprintf(f, `elem.On("click", h)`)
			fmt.Fprintln(f)
		}

		for _, e := range t.EventHandlers {
			fmt.Fprintf(f, `elem.%s(h)`, e.Name)
			fmt.Fprintln(f)
		}

		switch t.Type {
		case parent:
			fmt.Fprintln(f, `elem.Text("hello")`)
			fmt.Fprintln(f, `elem.Textf("hello %s", "Maxence")`)

		case privateParent:
			fmt.Fprintln(f, `elem.privateBody(Text("hello"))`)
		}

		fmt.Fprintln(f, "}")
	}
}
