package ticketswitch

import "github.com/shopspring/decimal"

// Discount contains all the information about the discount from the API
type Discount struct {
	AbsoluteSaving           decimal.Decimal `json:"absolute_saving"`
	AllowsLeavingSingleSeats string          `json:"allows_leaving_single_seats"`
	Code                     string          `json:"discount_code"`
	Description              string          `json:"discount_desc"`
	MinimumEligibleAge       int             `json:"discount_minimum_eligible_age"`
	MaximumEligibleAge       int             `json:"discount_maximum_eligible_age"`
	SemanticType             string          `json:"discount_semantic_type"`
	IsOffer                  bool            `json:"is_offer"`
	NonOfferSeatprice        decimal.Decimal `json:"non_offer_sale_seatprice"`
	NonOfferSurcharge        decimal.Decimal `json:"non_offer_sale_surcharge"`
	NonOfferCombined         decimal.Decimal `json:"non_offer_sale_combined"`
	NumberAvailable          int             `json:"number_available"`
	PercentageSaving         decimal.Decimal `json:"percentage_saving"`
	PriceBandCode            string          `json:"price_band_code"`
	Seatprice                decimal.Decimal `json:"sale_seatprice"`
	Surcharge                decimal.Decimal `json:"sale_surcharge"`
	Combined                 decimal.Decimal `json:"sale_combined"`
}

// DiscountsHolder is the intermediary discounts holder -- an artefact of the API
type DiscountsHolder struct {
	Discounts []Discount `json:"discount"`
}

// DiscountsResult  contains all the information from the GetDiscounts API call
type DiscountsResult struct {
	DiscountsHolder DiscountsHolder     `json:"discounts"`
	CurrencyCode    string              `json:"currency_code"`
	CurrencyDetails map[string]Currency `json:"currency_details"`
}
