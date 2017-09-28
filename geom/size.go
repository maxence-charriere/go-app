package geom

import "fmt"

// Size represents a two-dimensional size.
type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Area returns the area of the size.
func (s Size) Area() float64 {
	return s.Width * s.Height
}

func (s Size) String() string {
	return fmt.Sprintf("[%v * %v]", s.Width, s.Height)
}
