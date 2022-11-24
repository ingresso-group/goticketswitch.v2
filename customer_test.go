package ticketswitch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomer_Params(t *testing.T) {
	customer := Customer{
		AgentReference:             "1234567",
		FirstName:                  "Fred",
		LastName:                   "Flintstone",
		CountryCode:                "us",
		Title:                      "Mr.",
		Initials:                   "J.",
		Suffix:                     "CavM",
		AddressLineOne:             "1313 Boulder Lane",
		AddressLineTwo:             "The Quarry District",
		Postcode:                   "70777",
		Town:                       "Bedrock",
		County:                     "LA",
		EmailAddress:               "fred@rockslateandgravel.com",
		Phone:                      "001",
		WorkPhone:                  "002",
		HomePhone:                  "002",
		SupplierCanUseCustomerData: true,
		UserCanUseCustomerData:     true,
		WorldCanUseCustomerData:    true,
	}

	params := customer.Params()
	assert.Equal(t, "1234567", params["agent_ref"])
	assert.Equal(t, "Fred", params["first_name"])
	assert.Equal(t, "Flintstone", params["last_name"])
	assert.Equal(t, "us", params["country_code"])
	assert.Equal(t, "Mr.", params["title"])
	assert.Equal(t, "J.", params["initials"])
	assert.Equal(t, "CavM", params["suffix"])
	assert.Equal(t, "70777", params["postcode"])
	assert.Equal(t, "Bedrock", params["town"])
	assert.Equal(t, "LA", params["county"])
	assert.Equal(t, "fred@rockslateandgravel.com", params["email_address"])
	assert.Equal(t, "001", params["phone"])
	assert.Equal(t, "002", params["work_phone"])
	assert.Equal(t, "002", params["home_phone"])
	assert.Equal(t, "1", params["supplier_can_use_customer_data"])
	assert.Equal(t, "1", params["user_can_use_customer_data"])
	assert.Equal(t, "1", params["world_can_use_customer_data"])
	assert.Equal(t, "1313 Boulder Lane", params["address_line_one"])
	assert.Equal(t, "The Quarry District", params["address_line_two"])
	customer.SupplierCanUseCustomerData = false
	customer.UserCanUseCustomerData = false
	customer.WorldCanUseCustomerData = false

	params = customer.Params()
	assert.Equal(t, "0", params["supplier_can_use_customer_data"])
	assert.Equal(t, "0", params["user_can_use_customer_data"])
	assert.Equal(t, "0", params["world_can_use_customer_data"])
}
