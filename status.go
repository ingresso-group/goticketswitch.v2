package ticketswitch

import "time"

// Status describes the current state of a transaction.
type StatusResult struct {
	Languages        []string            `json:"language_list"`
	Trolley          Trolley             `json:"trolley_contents"`
	ReserveDatetime  time.Time           `json:"reserve_iso8601_date_and_time"`
	PurchaseDatetime time.Time           `json:"purchase_iso8601_date_and_time"`
	CurrencyDetails  map[string]Currency `json:"currency_details"`
	Status           string              `json:"transaction_status"`
}
