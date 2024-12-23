package resolv

import "math"

// Projection represents the projection of a shape (usually a ConvexPolygon) onto an axis for intersection testing.
// Normally, you wouldn't need to get this information, but it could be useful in some circumstances, I'm sure.
type Projection struct {
	Min, Max float64
}

// IsOverlapping returns whether a Projection is overlapping with the other, provided Projection. Credit to https://www.sevenson.com.au/programming/sat/
func (projection Projection) IsOverlapping(other Projection) bool {
	return projection.Overlap(other) > 0
}

// Overlap returns the amount that a Projection is overlapping with the other, provided Projection. Credit to https://dyn4j.org/2010/01/sat/#sat-nointer
func (projection Projection) Overlap(other Projection) float64 {
	return math.Min(projection.Max-other.Min, other.Max-projection.Min)
}

// IsInside returns whether the Projection is wholly inside of the other, provided Projection.
func (projection Projection) IsInside(other Projection) bool {
	return projection.Min >= other.Min && projection.Max <= other.Max
}
