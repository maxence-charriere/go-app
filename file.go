package app

// FilePanelConfig is a struct that describes a panel that selects files and
// directories to open.
type FilePanelConfig struct {
	// Reports whether the pannel allows multiple selection.
	MultipleSelection bool

	// Reports whether the pannel ignore directories.
	IgnoreDirectories bool

	// Reports whether the pannel ignore files.
	IgnoreFiles bool

	// Reports whether the pannel show hidden files.
	ShowHiddenFiles bool

	// Specify the file types to display in the pannel.
	// Accepts file extensions (eg. jpg, gif) and UTI (eg. public.jpeg).
	// Nil or empty allows all file types.
	FileTypes []string

	// If set, the function that is called when files or directories are
	// selected.
	OnSelect func(filenames []string) `json:"-"`
}

// SaveFilePanelConfig is a struct that describes a panel that selects a file to
// save.
type SaveFilePanelConfig struct {
	// Reports whether the pannel show hidden files.
	ShowHiddenFiles bool

	// Specify the file types to display in the pannel.
	// Accepts file extensions (eg. jpg, gif) and UTI (eg. public.jpeg).
	// Nil or empty allows all file types.
	FileTypes []string

	// If set, the function that is called when a file is selected
	OnSelect func(filename string) `json:"-"`
}
