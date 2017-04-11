package app

// FilePicker is a struct that describes a file picker.
// It will be used by a driver to create a native file picker that allow to
// select files and directories filenames.
type FilePicker struct {
	MultipleSelection bool
	NoDir             bool
	NoFile            bool
	OnPick            func(filenames []string)
}

// NewFilePicker creates and opens a native file picker described by fp.
func NewFilePicker(p FilePicker) Elementer {
	return driver.NewElement(p)
}
