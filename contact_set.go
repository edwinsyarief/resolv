package resolv

import (
	"math"

	ebimath "github.com/edwinsyarief/ebi-math"
)

// Intersection represents a single point of contact against a line or surface.
type Intersection struct {
	Point  ebimath.Vector // The point of contact.
	Normal ebimath.Vector // The normal of the surface contacted.
}

// IntersectionSet represents a set of intersections between the calling object and one other intersecting Shape.
// A Shape's intersection test may iterate through multiple IntersectionSets - one for each pair of intersecting objects.
type IntersectionSet struct {
	Intersections []Intersection // Slice of points indicating contact between the two Shapes.
	Center        ebimath.Vector // Center of the Contact set; this is the average of all Points contained within all contacts in the IntersectionSet.
	MTV           ebimath.Vector // Minimum Translation ebimath.Vector; this is the ebimath.Vector to move a Shape on to move it to contact with the other, intersecting / contacting Shape.
	OtherShape    IShape         // The other shape involved in the contact.
}

func newIntersectionSet() IntersectionSet {
	return IntersectionSet{}
}

// LeftmostPoint returns the left-most point out of the IntersectionSet's Points slice.
// If the IntersectionSet is empty, this returns a zero ebimath.Vector.
func (is IntersectionSet) LeftmostPoint() ebimath.Vector {

	var left ebimath.Vector
	set := false

	for _, contact := range is.Intersections {

		if !set || contact.Point.X < left.X {
			left = contact.Point
			set = true
		}

	}

	return left

}

// RightmostPoint returns the right-most point out of the IntersectionSet's Points slice.
// If the IntersectionSet is empty, this returns a zero ebimath.Vector.
func (is IntersectionSet) RightmostPoint() ebimath.Vector {

	var right ebimath.Vector
	set := false

	for _, contact := range is.Intersections {

		if !set || contact.Point.X > right.X {
			right = contact.Point
			set = true
		}

	}

	return right

}

// TopmostPoint returns the top-most point out of the IntersectionSet's Points slice. I
// f the IntersectionSet is empty, this returns a zero ebimath.Vector.
func (is IntersectionSet) TopmostPoint() ebimath.Vector {

	var top ebimath.Vector
	set := false

	for _, contact := range is.Intersections {

		if !set || contact.Point.Y < top.Y {
			top = contact.Point
			set = true
		}

	}

	return top

}

// BottommostPoint returns the bottom-most point out of the IntersectionSet's Points slice.
// If the IntersectionSet is empty, this returns a zero ebimath.Vector.
func (is IntersectionSet) BottommostPoint() ebimath.Vector {

	var bottom ebimath.Vector
	set := false

	for _, contact := range is.Intersections {

		if !set || contact.Point.Y > bottom.Y {
			bottom = contact.Point
			set = true
		}

	}

	return bottom

}

// IsEmpty returns if the IntersectionSet is empty (and so contains no points of iontersection). This should never actually be true.
func (is IntersectionSet) IsEmpty() bool {
	return len(is.Intersections) == 0
}

// Distance returns the distance between all of the intersection points when projected against an axis.
func (is IntersectionSet) Distance(alongAxis ebimath.Vector) float64 {
	alongAxis = alongAxis.Unit()
	top, bottom := math.MaxFloat64, -math.MaxFloat64
	for _, c := range is.Intersections {
		d := alongAxis.Dot(c.Point)
		top = min(top, d)
		bottom = max(bottom, d)
	}
	return bottom - top
}
