package app

// Docker is the interface that describes a dock.
type Docker interface {
	Contexter

	// SetIcon set the dock icon to the image targeted by imageName.
	SetIcon(imageName string)

	// SetBadge set the dock badge with the string value of v.
	SetBadge(v interface{})
}
