package ticketswitch

type Source struct {
	// The code for the source
	Code string `json:"source_code"`
	// The description of the source
	Description string `json:"source_desc_from_config"`
	// The after sales email of the source
	Email string `json:"source_after_sales_email"`
	// Postal address of the source
	Address string `json:"source_postal_addr"`
	// Class of the source
	Class string `json:"source_system_class"`
	// Type of the string
	Type string `json:"source_system_type_string"`
	// Terms and Conditions
	TermsAndConditions string `json:"source_t_and_c"`
}

type SourcesResult struct {
	Sources []Source
}
