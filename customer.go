package ticketswitch

// Customer contains information about the customer that bought tickets
type Customer struct {
	AgentReference             string `json:"agent_ref"`
	FirstName                  string `json:"first_name"`
	LastName                   string `json:"last_name"`
	CountryCode                string `json:"country_code"`
	Title                      string `json:"title"`
	Initials                   string `json:"initials"`
	Suffix                     string `json:"suffix"`
	Postcode                   string `json:"postcode"`
	Town                       string `json:"town"`
	County                     string `json:"county"`
	EmailAddress               string `json:"email_addr"`
	Phone                      string `json:"phone"`
	WorkPhone                  string `json:"work_phone"`
	HomePhone                  string `json:"home_phone"`
	AddressLineOne             string `json:"addr_line_one"`
	AddressLineTwo             string `json:"addr_line_two"`
	SupplierCanUseCustomerData bool   `json:"supplier_can_use_customer_data"`
	UserCanUseCustomerData     bool   `json:"user_can_use_customer_data"`
	WorldCanUseCustomerData    bool   `json:"world_can_use_customer_data"`
}

// Params returns the customer data as a map
func (customer *Customer) Params() map[string]string {
	values := map[string]string{
		"agent_ref":                      customer.AgentReference,
		"first_name":                     customer.FirstName,
		"last_name":                      customer.LastName,
		"country_code":                   customer.CountryCode,
		"title":                          customer.Title,
		"initials":                       customer.Initials,
		"suffix":                         customer.Suffix,
		"postcode":                       customer.Postcode,
		"town":                           customer.Town,
		"county":                         customer.County,
		"email_address":                  customer.EmailAddress,
		"phone":                          customer.Phone,
		"work_phone":                     customer.WorkPhone,
		"home_phone":                     customer.HomePhone,
		"address_line_one":               customer.AddressLineOne,
		"address_line_two":               customer.AddressLineTwo,
		"supplier_can_use_customer_data": "0",
		"user_can_use_customer_data":     "0",
		"world_can_use_customer_data":    "0",
	}

	if customer.SupplierCanUseCustomerData {
		values["supplier_can_use_customer_data"] = "1"
	}
	if customer.UserCanUseCustomerData {
		values["user_can_use_customer_data"] = "1"
	}
	if customer.WorldCanUseCustomerData {
		values["world_can_use_customer_data"] = "1"
	}

	return values
}
