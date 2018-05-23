package ticketswitch

import (
	"testing"

	geo "github.com/kellydunn/golang-geo"
	"github.com/stretchr/testify/assert"
)

func TestCircle(t *testing.T) {
	area := NewCircle(45.67890, 98.76543, 123.4567)
	assert.Equal(t, 45.67890, area.Lat())
	assert.Equal(t, 98.76543, area.Lng())
	assert.Equal(t, 123.4567, area.Radius)

	point := geo.NewPoint(12.345, 67.890)
	area = NewCircleWithPoint(*point, 987.65)
	assert.Equal(t, 12.345, area.Lat())
	assert.Equal(t, 67.890, area.Lng())
	assert.Equal(t, 987.65, area.Radius)

}

func TestCircle_Valid(t *testing.T) {
	area := NewCircle(45.67890, 98.76543, 123.4567)
	assert.True(t, area.Valid())

	area = NewCircle(91.0, 98.76543, 123.4567)
	assert.False(t, area.Valid())

	area = NewCircle(-91.0, 98.76543, 123.4567)
	assert.False(t, area.Valid())

	area = NewCircle(45.67890, 181.0, 123.4567)
	assert.False(t, area.Valid())

	area = NewCircle(45.67890, -181.0, 123.4567)
	assert.False(t, area.Valid())

	area = NewCircle(45.67890, 98.76543, 0)
	assert.False(t, area.Valid())

	area = NewCircle(45.67890, 98.76543, -1.0)
	assert.False(t, area.Valid())
}

func TestCircle_Param(t *testing.T) {
	area := &Circle{}
	assert.Equal(t, "0.000000:0.000000:0.000000", area.Param())

	area = NewCircle(45.67890, 98.76543, 123.4567)
	assert.Equal(t, "45.678900:98.765430:123.456700", area.Param())
}
