package ticketswitch

import (
	"time"

	"github.com/shopspring/decimal"
)

// AvailabilityDetails describes an availability summary for a performance.
// This information is generated using cached data.
type AvailabilityDetails struct {
	// identifier of the ticket type.
	TicketTypeCode string
	// human readable description of the ticket type.
	TicketTypeDescription string
	// identifier of the price band.
	PriceBandCode string
	// human readable description of the price band.
	PriceBandDescription string
	// price of an individual seat.
	SeatPrice decimal.Decimal
	// additional charges per seat.
	Surcharge decimal.Decimal
	// the non-offer price of an individual seat.
	FullSeatPrice decimal.Decimal
	// the non-offer additional charges per seat.
	FullSurcharge decimal.Decimal
	// the currency of the prices.
	CurrencyCode string
	// the first date and time this combination of ticket type and price
	// band is available from.
	FirstDate time.Time
	// the latest date and time this combination of ticket type and price
	// band is available from.
	LastDate time.Time
	// list of valid number of tickets available for selection.
	ValidQuantities []int
}
