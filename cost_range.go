package ticketswitch

import "github.com/shopspring/decimal"

// CostRange gives summarized pricing for events and performances.
//
// This information is returned from cached data collected when making actual
// calls to the backend system, and should not be considered accurate.
type CostRange struct {
	// list of valid quanities available for purchase.
	ValidQuanities []int `json:"valid_quanities"`

	// the minimum cost per seat the customer might be expected to pay.
	MinSeatPrice decimal.Decimal `json:"min_seatprice"`
	// the maximum cost per seat the customer might be expected to pay.
	MaxSeatPrice decimal.Decimal `json:"max_seatprice"`
	// the minimum surcharge per seat the customer might be expected to pay.
	MinSurcharge decimal.Decimal `json:"min_surcharge"`
	// the maximum surcharge per seat the customer might be expected to pay.
	MaxSurcharge decimal.Decimal `json:"max_surcharge"`

	// currency the cost range and offer prices are in.
	CurrencyCode string   `json:"currency_code"`
	Currency     Currency `json:"currency"`

	// offer with the highest percentage saving.
	BestValueOffer Offer `json:"best_value_offer"`

	// offer with the highest absolute saving.
	MaxSavingOffer Offer `json:"max_saving_offer"`

	// offer with the lowest cost.
	MinCostOffer Offer `json:"min_cost_offer"`

	//  offer with the top price.
	TopPrice Offer `json:"top_price"`
}

// CostRangeDetails Summarizes pricing by ticket types and price bands for an
// event/performance
//
// This information is returned from cached data collected when making actual
// calls to the backend system, and should not be considered accurate.
type CostRangeDetails struct {
	TicketTypeCode     string    `json:"ticket_type_code"`
	PriceBandCode      string    `json:"price_band_code"`
	TicketTypeDesc     string    `json:"ticket_type_desc"`
	PriceBandDesc      string    `json:"price_band_desc"`
	CostRange          CostRange `json:"cost_range"`
	NoSinglesCostRange CostRange `json:"no_singles_cost_range"`
}
