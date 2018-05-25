package ticketswitch

import "github.com/shopspring/decimal"

type SendMethod struct {
	Code string          `json:"send_code"`
	Cost decimal.Decimal `json:"send_cost"`
	Desc string          `json:"send_desc"`
	Type string          `json:"send_type"`
}
