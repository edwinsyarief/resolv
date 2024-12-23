package resolv

import (
	"sort"

	ebimath "github.com/edwinsyarief/ebi-math"
)

// Set represents a Set of elements.
type Set[E comparable] map[E]struct{}

// newSet creates a new set.
func newSet[E comparable]() Set[E] {
	return Set[E]{}
}

// Clone clones the Set.
func (s Set[E]) Clone() Set[E] {
	newSet := newSet[E]()
	newSet.Combine(s)
	return newSet
}

// Set sets the Set to have the same values as in the given other Set.
func (s Set[E]) Set(other Set[E]) {
	s.Clear()
	s.Combine(other)
}

// Add adds the given elements to a set.
func (s Set[E]) Add(elements ...E) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

// Combine combines the given other elements to the set.
func (s Set[E]) Combine(otherSet Set[E]) {
	for element := range otherSet {
		s.Add(element)
	}
}

// Contains returns if the set contains the given element.
func (s Set[E]) Contains(element E) bool {
	_, ok := s[element]
	return ok
}

// Remove removes the given element from the set.
func (s Set[E]) Remove(elements ...E) {
	for _, element := range elements {
		delete(s, element)
	}
}

// Clear clears the set.
func (s Set[E]) Clear() {
	for v := range s {
		delete(s, v)
	}
}

// ForEach runs the provided function for each element in the set.
func (s Set[E]) ForEach(f func(element E) bool) {
	for element := range s {
		if !f(element) {
			break
		}
	}
}

/////

// shapeIDSet is an easy way to determine if a shape has been iterated over before (used for filtering through shapes from CellSelections).
type shapeIDSet []uint32

func (s shapeIDSet) idInSet(id uint32) bool {
	for _, v := range s {
		if v == id {
			return true
		}
	}
	return false
}

var cellSelectionForEachIDSet = shapeIDSet{}

// LineTestSettings is a struct of settings to be used when performing line tests (the equivalent of 3D hitscan ray tests for 2D)
type LineTestSettings struct {
	Start       ebimath.Vector // The start of the line to test shapes against
	End         ebimath.Vector // The end of the line to test chapes against
	TestAgainst ShapeIterator  // The collection of shapes to test against
	// The callback to be called for each intersection between the given line, ranging from start to end, and each shape given in TestAgainst.
	// set is the intersection set that contains information about the intersection, index is the index of the current index
	// and count is the total number of intersections detected from the intersection test.
	// The boolean the callback returns indicates whether the LineTest function should continue testing or stop at the currently found intersection.
	OnIntersect  func(set IntersectionSet, index, max int) bool
	callingShape IShape
}

var intersectionSets []IntersectionSet

// LineTest instantly tests a selection of shapes against a ray / line.
// Note that there is no MTV for these results.
func LineTest(settings LineTestSettings) bool {

	castMargin := 0.01 // Basically, the line cast starts are a smidge back so that moving to contact doesn't make future line casts fail
	vu := settings.End.Sub(settings.Start).Unit()
	start := settings.Start.Sub(vu.ScaleF(castMargin))

	line := newCollidingLine(start.X, start.Y, settings.End.X, settings.End.Y)

	intersectionSets = intersectionSets[:0]

	i := 0

	settings.TestAgainst.ForEach(func(other IShape) bool {

		if other == settings.callingShape {
			return true
		}

		i++

		contactSet := newIntersectionSet()

		switch shape := other.(type) {

		case *Circle:

			res := line.IntersectionPointsCircle(shape)

			if len(res) > 0 {
				for _, contactPoint := range res {
					contactSet.Intersections = append(contactSet.Intersections, Intersection{
						Point:  contactPoint,
						Normal: contactPoint.Sub(shape.position).Unit(),
					})
				}
			}

		case *ConvexPolygon:

			for _, otherLine := range shape.Lines() {

				if point, ok := line.IntersectionPointsLine(otherLine); ok {
					contactSet.Intersections = append(contactSet.Intersections, Intersection{
						Point:  point,
						Normal: otherLine.Normal(),
					})
				}

			}

		}

		if len(contactSet.Intersections) > 0 {

			contactSet.OtherShape = other

			for _, contact := range contactSet.Intersections {
				contactSet.Center = contactSet.Center.Add(contact.Point)
			}

			contactSet.Center.X /= float64(len(contactSet.Intersections))
			contactSet.Center.Y /= float64(len(contactSet.Intersections))

			// Sort the points by distance to line start
			sort.Slice(contactSet.Intersections, func(i, j int) bool {
				return contactSet.Intersections[i].Point.DistanceSquaredTo(settings.Start) < contactSet.Intersections[j].Point.DistanceSquaredTo(settings.Start)
			})

			contactSet.MTV = contactSet.Intersections[0].Point.Sub(settings.Start).Sub(vu.ScaleF(castMargin))

			intersectionSets = append(intersectionSets, contactSet)

		}

		return true

	})

	// Sort intersection sets by distance from closest hit to line start
	sort.Slice(intersectionSets, func(i, j int) bool {
		return intersectionSets[i].Intersections[0].Point.DistanceSquaredTo(line.Start) < intersectionSets[j].Intersections[0].Point.DistanceSquaredTo(line.Start)
	})

	// Loop through all intersections and iterate through them
	if settings.OnIntersect != nil {
		for i, c := range intersectionSets {
			if !settings.OnIntersect(c, i, len(intersectionSets)) {
				break
			}
		}
	}

	return len(intersectionSets) > 0

}
