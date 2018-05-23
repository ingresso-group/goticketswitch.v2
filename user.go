package ticketswitch

// User describes a user of the API
type User struct {
	// the user identifier.
	ID string `json:"user_id"`
	// human readable name.
	Name string `json:"real_name"`
	// ISO 3166-1 country code.
	Country string `json:"default_country_code"`
	// the identifier of the sub user.
	SubUser string `json:"sub_user"`
	// indicates that the account is a b2b account.
	IsB2B bool `json:"is_b2b"`
	// what will appear on a customers bank statement when ingresso takes the
	// payment.
	StatementDescriptor string `json:"statement_descriptor"`
	// what product group a user belongs to. Users in the same backend_group
	// will see the same products.
	BackendGroup string `json:"backend_group"`
	// what content group a user belongs to. Users in the same content_group
	// will see the same textual and graphical content.
	ContentGroup string `json:"content_group"`
}
