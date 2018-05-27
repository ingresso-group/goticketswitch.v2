package ticketswitch

import (
	"time"

	"github.com/shopspring/decimal"
)

// PaymentMethod defines an interface that can be used for supplying payment
// parameters to the API
type PaymentMethod interface {
	PaymentParams() map[string]string
}

// Debitor represents information about a 3rd party that will take payment from
// your customer.
//
// This information is primarily used for bypassing callouts and integrating
// directly with payment providers on the front end.
//
// When your account is set up to sell on credit (i.e. you are always taking
// payment from the customer in your application directly), this information
// will not be present, and it should not be relevant.
//
// When the source system is taking payment this information will not be
// present.
//
// When debitor information is not present or you are not front end
// integrating you should refer to the Reservation.NeedsPaymentCard,
// Reservation.NeedsEmailAddress, and Reservation.NeedsAgentReference
// as to what information you need to pass back to the API for purchasing
// tickets.
//
// Regardless of the debitor it's advisable to implement the full
// purchase/callout/callback process in the event that your front end
// integration goes awry.
type Debitor struct {
	Type            string                 `json:"debitor_type"`
	Name            string                 `json:"debitor_name"`
	Description     string                 `json:"debitor_desc"`
	IntegrationData map[string]interface{} `json:"debitor_integration_data"`
	AggregrationKey string                 `json:"debitor_aggregration_key"`
}

// Callout describes how a customer should be redirected in order to provide
// additional data, for example for 3D secure, or logging into paypal.
type Callout struct {
	Code            string                 `json:"bundle_source_code"`
	Description     string                 `json:"bundle_source_desc"`
	Total           decimal.Decimal        `json:"bundle_total_cost"`
	Type            string                 `json:"callout_type"`
	Destination     string                 `json:"callout_destination_url"`
	Parameters      map[string]string      `json:"callout_parameters"`
	IntegrationData map[string]interface{} `json:"callout_integration_data"`
	Debitor         Debitor                `json:"debitor"`
	CurrencyCode    string                 `json:"currency_code"`
	ReturnToken     string                 `json:"return_token"`
}

// MakePurchaseResult is the result from the MakePurchase client call.
type MakePurchaseResult struct {
	Status           string              `json:"transaction_status"`
	Callout          *Callout            `json:"callout,omitempty"`
	PendingCallout   *Callout            `json:"pending_callout,omitempty"`
	Currency         map[string]Currency `json:"currency_details"`
	Trolley          Trolley             `json:"trolley_contents"`
	Customer         Customer            `json:"customer"`
	ReserveDatetime  time.Time           `json:"reserve_iso8601_date_and_time"`
	PurchaseDatetime time.Time           `json:"purchase_iso8601_date_and_time"`
	ReserveUser      User                `json:"reserve_user"`
	Languages        []string            `json:"language_list"`
}
