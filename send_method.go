package ticketswitch

import "github.com/shopspring/decimal"

type Country struct {
	Code string `json:"country_code"`
	Desc string `json:"country_desc"`
}

type PermittedCountries struct {
	Countries []Country `json:"country"`
}

type SendMethod struct {
	Code                 string             `json:"send_code"`
	Cost                 decimal.Decimal    `json:"send_cost"`
	Desc                 string             `json:"send_desc"`
	Type                 string             `json:"send_type"`
	PermittedCountries   PermittedCountries `json:"permitted_countries"`
	FinalType            string             `json:"send_final_tpe"`
	CanGenerateSelfPrint bool               `json:"CanGenerateSelfPrint"`
	SelfPrintVoucherURL  string             `json:"self_print_voucher_url"`
	HasHTMLPage          bool               `json:"has_html_page"`
}

type SendMethodsHolder struct {
	SendMethods []SendMethod `json:"send_method"`
}

type SendMethodsResults struct {
	CurrencyDetails   map[string]Currency `json:"currency_details"`
	CurrencyCode      string              `json:"currency_code"`
	SourceCode        string              `json:"source_code"`
	SendMethodsHolder SendMethodsHolder   `json:"send_methods"`
}
