package ticketswitch

// Currency contains the information about a currency.
type Currency struct {
	//  ISO 4217 currency code.
	Code string `json:"currency_code"`
	// precision of decimal numbers.
	Places int `json:"currency_places"`
	// a symbol to place before the digits of a price.
	PreSymbol string `json:"currency_pre_symbol"`
	// a symbol to place after the digits of a price.
	PostSymbol string `json:"currency_post_symbol"`
	// arbitrary scale factor, can be used to roughly convert from one currency
	// to another.
	Factor int `json:"currency_factor"`
	//  internal identifier.
	Number int `json:"currency_number"`
}
