package ticketswitch

import "github.com/shopspring/decimal"

type Offer struct {
	// the price per seat/ticket.
	SeatPrice decimal.Decimal `json:"offer_seatprice"`
	// the additional charges per seat/ticket.
	Surcharge decimal.Decimal `json:"offer_surcharge"`
	// the original price per seat/ticket.
	FullSeatPrice decimal.Decimal `json:"full_seatprice"`
	// the original additional charges per seat/ticket.
	FullSurcharge decimal.Decimal `json:"full_surcharge"`
	// the amount of money saved by this offer.
	AbsoluteSaving decimal.Decimal `json:"absolute_saving"`
	// the amount of money saved by this offer, as a percentage of the original
	// price.
	PercentageSaving decimal.Decimal `json:"percentage_saving"`
}
