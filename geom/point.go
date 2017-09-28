package geom

import (
	"fmt"
)

// Point represents a Cartesian coordinate.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (p Point) String() string {
	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}
