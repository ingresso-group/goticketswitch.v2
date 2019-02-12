package ticketswitch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateRange(t *testing.T) {
	assert.Equal(t, "", DateRange(time.Time{}, time.Time{}))
	from := time.Date(2017, 8, 7, 12, 21, 25, 0, time.UTC)
	to := time.Date(2018, 11, 21, 19, 30, 0, 0, time.UTC)

	assert.Equal(t, "20170807:", DateRange(from, time.Time{}))
	assert.Equal(t, ":20181121", DateRange(time.Time{}, to))
	assert.Equal(t, "20170807:20181121", DateRange(from, to))
}
