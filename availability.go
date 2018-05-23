package ticketswitch

import (
	"strconv"

	"github.com/shopspring/decimal"
)

// PriceBand describes a subset of available tickets within a ticket type
// defined by price point. The price of a price band is defined by it's default
// discount code, this is normally the most expensive discount option available
type PriceBand struct {
	Code                     string                `json:"price_band_code"`
	DiscountCode             string                `json:"discount_code"`
	DiscountDesc             string                `json:"discount_desc"`
	NumberAvailable          int                   `json:"number_available"`
	Seatprice                decimal.Decimal       `json:"sale_seatprice"`
	Surcharge                decimal.Decimal       `json:"sale_surcharge"`
	AllowsLeavingSingleSeats string                `json:"allows_leaving_single_seats"`
	IsOffer                  bool                  `json:"is_offer"`
	NonOfferSeatprice        decimal.Decimal       `json:"non_offer_sale_seatprice"`
	NonOfferSurcharge        decimal.Decimal       `json:"non_offer_sale_surcharge"`
	PercentageSaving         decimal.Decimal       `json:"percentage_saving"`
	AbsoluteSaving           decimal.Decimal       `json:"absolute_saving"`
	FreeSeatBlocksRaw        map[string][][]string `json:"free_seat_blocks"`
	RestrictedViewSeatsRaw   []string              `json:"restricted_view_seats_raw"`
	SeatsByTextMessageRaw    []string              `json:"seats_by_text_message_raw"`
	PredictedUserCommission  UserCommission        `json:"predicted_user_commission"`
}

// TicketType describes a sub set of available tickets defined by some non
// price related parameters. Normally for venue based performances this will
// indicate a part of house or area within the venue.
type TicketType struct {
	Code       string      `json:"ticket_type_code"`
	Desc       string      `json:"ticket_type_desc"`
	PriceBands []PriceBand `json:"price_band"`
}

// Availability holds the details of the availble tickets
type Availability struct {
	TicketTypes []TicketType `json:"ticket_type"`
}

// AvailabilityResult describes the current state of available seats for a
// Performance
type AvailabilityResult struct {
	Availability                Availability        `json:"availability"`
	BackendIsBroken             bool                `json:"backend_is_broken"`
	BackendIsDown               bool                `json:"backend_is_down"`
	BackendThrottleFailed       bool                `json:"backend_throttle_failed"`
	ContiguousSeatSelectionOnly bool                `json:"contiguous_seat_selection_only"`
	CurrencyCode                string              `json:"currency_code"`
	CurrencyDetails             map[string]Currency `json:"currency_details"`
	ValidQuantities             []int               `json:"valid_quantities"`
}

// GetAvailabilityParams are parameters that can be passed to the
// GetAvailability call.
type GetAvailabilityParams struct {
	UniversalParams
	NumberOfSeats  int
	Discounts      bool
	ExampleSeats   bool
	SeatBlocks     bool
	UserCommission bool
}

// Params returns the call parameters as a map
func (params *GetAvailabilityParams) Params() map[string]string {
	values := make(map[string]string)

	if params.NumberOfSeats > 0 {
		values["number_of_seats"] = strconv.Itoa(params.NumberOfSeats)
	}

	if params.Discounts {
		values["add_discounts"] = "1"
	}

	if params.ExampleSeats {
		values["add_example_seats"] = "1"
	}

	if params.SeatBlocks {
		values["add_seat_blocks"] = "1"
	}

	if params.UserCommission {
		values["req_predicted_commission"] = "1"
	}

	for k, v := range params.Universal() {
		values[k] = v
	}

	return values
}
