package ticketswitch

import (
	"fmt"

	geo "github.com/kellydunn/golang-geo"
)

// Circle represents a geographical point and an area around it.
type Circle struct {
	geo.Point
	Radius float64
}

// NewCircle returns a pointer to a newly created GeoArea
func NewCircle(lat, long, radius float64) *Circle {
	point := geo.NewPoint(lat, long)
	area := &Circle{
		Point:  *point,
		Radius: radius,
	}
	return area
}

// NewCircleWithPoint returns a pointer to a newly created GeoArea using an
// existing geo.Point
func NewCircleWithPoint(point geo.Point, radius float64) *Circle {
	area := &Circle{
		Point:  point,
		Radius: radius,
	}
	return area
}

// Valid checks if the Circle is valid
func (area *Circle) Valid() bool {
	if area.Radius <= 0 {
		return false
	}

	if area.Lat() < -90.0 || area.Lat() > 90.0 {
		return false
	}

	if area.Lng() < -180.0 || area.Lng() > 180.0 {
		return false
	}

	return true
}

// Param returns a string representation of a set of coordinates and a radius for
// consumption by the API.
func (area *Circle) Param() string {
	return fmt.Sprintf("%.6f:%.6f:%.6f", area.Lat(), area.Lng(), area.Radius)
}
