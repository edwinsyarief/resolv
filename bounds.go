package resolv

import (
	"math"

	ebimath "github.com/edwinsyarief/ebi-math"
)

// Bounds represents the minimum and maximum bounds of a Shape.
type Bounds struct {
	Min, Max ebimath.Vector
	space    *Space
}

func (b Bounds) toCellSpace() (int, int, int, int) {

	minX := int(math.Floor(b.Min.X / float64(b.space.cellWidth)))
	minY := int(math.Floor(b.Min.Y / float64(b.space.cellHeight)))
	maxX := int(math.Floor(b.Max.X / float64(b.space.cellWidth)))
	maxY := int(math.Floor(b.Max.Y / float64(b.space.cellHeight)))

	return minX, minY, maxX, maxY
}

// Center returns the center position of the Bounds.
func (b Bounds) Center() ebimath.Vector {
	return b.Min.Add(b.Max.Sub(b.Min).ScaleF(0.5))
}

// Width returns the width of the Bounds.
func (b Bounds) Width() float64 {
	return b.Max.X - b.Min.X
}

// Height returns the height of the bounds.
func (b Bounds) Height() float64 {
	return b.Max.Y - b.Min.Y
}

// Move moves the Bounds, such that the center point is offset by {x, y}.
func (b Bounds) Move(x, y float64) Bounds {
	b.Min.X += x
	b.Min.Y += y
	b.Max.X += x
	b.Max.Y += y
	return b
}

// MoveVec moves the Bounds by the vector provided, such that the center point is offset by {x, y}.
func (b *Bounds) MoveVec(vec ebimath.Vector) Bounds {
	return b.Move(vec.X, vec.Y)
}

// IsIntersecting returns if the Bounds is intersecting with the given other Bounds.
func (b Bounds) IsIntersecting(other Bounds) bool {
	bounds := b.Intersection(other)
	return !bounds.IsEmpty()
}

// Intersection returns the intersection between the two Bounds objects.
func (b Bounds) Intersection(other Bounds) Bounds {

	overlap := Bounds{}

	if other.Max.X < b.Min.X || other.Min.X > b.Max.X || other.Max.Y < b.Min.Y || other.Min.Y > b.Max.Y {
		return overlap
	}

	overlap.Min.X = math.Min(other.Max.X, b.Max.X)
	overlap.Max.X = math.Max(other.Min.X, b.Min.X)

	overlap.Min.Y = math.Min(other.Max.Y, b.Max.Y)
	overlap.Max.Y = math.Max(other.Min.Y, b.Min.Y)

	return overlap

}

// IsEmpty returns true if the Bounds's minimum and maximum corners are 0.
func (b Bounds) IsEmpty() bool {
	return b.Max.X-b.Min.X == 0 && b.Max.Y-b.Min.X == 0
}
