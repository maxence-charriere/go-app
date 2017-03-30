package app

// FileChooser represents a panel that allow to select files and directories.
type FileChooser struct {
	MultipleSelection bool
	NoDir             bool
	NoFile            bool
	OnChoose          func(filenames []string)
}

// OpenFileChooser opens fc.
func NewFileChooser(fc FileChooser) Elementer {
	return driver.NewElement(fc)
}
