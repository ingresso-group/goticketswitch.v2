package ticketswitch

import "github.com/shopspring/decimal"

type SendMethod struct {
	Code                 string          `json:"send_code"`
	Cost                 decimal.Decimal `json:"send_cost"`
	Desc                 string          `json:"send_desc"`
	Type                 string          `json:"send_type"`
	FinalType            string          `json:"send_final_tpe"`
	CanGenerateSelfPrint bool            `json:"CanGenerateSelfPrint"`
	SelfPrintVoucherURL  string          `json:"self_print_voucher_url"`
	HasHTMLPage          bool            `json:"has_html_page"`
}
