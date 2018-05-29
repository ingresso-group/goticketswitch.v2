package ticketswitch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMakeReservationParams_Params(t *testing.T) {
	// Minimum parameter set
	params := MakeReservationParams{
		NumberOfSeats:  3,
		PerformanceID:  "6IF-C30",
		PriceBandCode:  "C/pool",
		TicketTypeCode: "CIRCLE",
	}

	values := params.Params()
	assert.Equal(t, values["no_of_seats"], "3")
	assert.Equal(t, values["perf_id"], "6IF-C30")
	assert.Equal(t, values["price_band_code"], "C/pool")
	assert.Equal(t, values["ticket_type_code"], "CIRCLE")

	// Test with all the extras:
	params.DepartureDate = time.Date(2018, 5, 24, 15, 13, 7, 0, time.UTC)
	params.Discounts = []string{"CHILD", "ADULT", "ADULT"}
	params.Seats = []string{"A1", "A2", "A3"}
	params.SendMethod = "POST"
	params.SourceCode = "ext_test0"
	params.TrolleyToken = "alkdfja8sldifa9oiefjaeiojfa2eijfasekfjasldkfjasdasdlkfhaskduhfaksjudf"
	params.UserCommission = true

	values = params.Params()
	assert.Equal(t, values["departure_date"], "20180524")
	assert.Equal(t, values["disc0"], "CHILD")
	assert.Equal(t, values["disc1"], "ADULT")
	assert.Equal(t, values["disc2"], "ADULT")
	assert.Equal(t, values["seat0"], "A1")
	assert.Equal(t, values["seat1"], "A2")
	assert.Equal(t, values["seat2"], "A3")
	assert.Equal(t, values["ext_test0_send_code"], "POST")
	assert.Equal(t, values["trolley_token"], "alkdfja8sldifa9oiefjaeiojfa2eijfasekfjasldkfjasdasdlkfhaskduhfaksjudf")
	assert.Equal(t, values["req_predicted_commission"], "1")
}
