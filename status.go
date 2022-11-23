package ticketswitch

import "time"

// StatusResult describes the current state of a transaction.
type StatusResult struct {
    Languages            []string            `json:"language_list"`
    MinutesLeftOnReserve float64             `json:"minutes_left_on_reserve"`
    Trolley              Trolley             `json:"trolley_contents"`
    ReserveDatetime      time.Time           `json:"reserve_iso8601_date_and_time"`
    PurchaseDatetime     time.Time           `json:"purchase_iso8601_date_and_time"`
    CurrencyDetails      map[string]Currency `json:"currency_details"`
    Customer             Customer            `json:"customer"`
    Status               string              `json:"transaction_status"`
}
