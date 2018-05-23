package ticketswitch

import "github.com/shopspring/decimal"

// UserCommission describes how much a user will be paid for selling a ticket
type UserCommission struct {
	IncVat       decimal.Decimal `json:"amount_excluding_vat"`
	ExVat        decimal.Decimal `json:"amount_including_vat"`
	CurrencyCode string          `json:"commission_currency_code"`
}
