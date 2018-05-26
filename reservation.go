package ticketswitch

import (
	"fmt"
	"strconv"
	"time"
)

// ReservationResult contains the results of MakeReservation call
type ReservationResult struct {
	AllowedCountries               map[string]string   `json:"allowed_countries"`
	CanEditAddress                 bool                `json:"can_edit_address"`
	CurrencyDetails                map[string]Currency `json:"currency_details"`
	InputContainedUnavailableOrder bool                `json:"input_contained_unavailable_order"`
	LanguageList                   []string            `json:"language_list"`
	MinutesLeftOnReserve           float64             `json:"minutes_left_on_reserve"`
	NeedsAgentReference            bool                `json:"needs_agent_reference"`
	NeedsEmailAddress              bool                `json:"needs_agent_reference"`
	NeedsPaymentCard               bool                `json:"needs_agent_reference"`
	PrefilledAddress               map[string]string   `json:"prefilled_address"`
	ReserveTime                    time.Time           `json:"reserve_iso8601_date_and_time"`
	TransactionStatus              string              `json:"transaction_status"`
	Trolley                        Trolley             `json:"trolley_contents"`
	UnreservedOrders               []Order             `json:"unreserved_orders"`
}

// MakeReservationParams stores the parameters for making a reservation
type MakeReservationParams struct {
	UniversalParams
	DepartureDate  time.Time
	Discounts      []string
	NumberOfSeats  int    // Required
	PerformanceID  string // Required
	PriceBandCode  string // Required
	Seats          []string
	SendMethod     string
	SourceCode     string // Required if specifying the send method
	TicketTypeCode string // Required
	TrolleyToken   string
	UserCommission bool
}

// Params returns the call parameters as a map
func (params *MakeReservationParams) Params() map[string]string {
	values := make(map[string]string)

	values["no_of_seats"] = strconv.Itoa(params.NumberOfSeats)
	values["perf_id"] = params.PerformanceID
	values["price_band_code"] = params.PriceBandCode
	values["ticket_type_code"] = params.TicketTypeCode

	if !params.DepartureDate.IsZero() {
		values["departure_date"] = params.DepartureDate.Format("20060102")
	}

	for index, disc := range params.Discounts {
		values[fmt.Sprintf("disc%d", index)] = disc
	}

	for index, seat := range params.Seats {
		values[fmt.Sprintf("seat%d", index)] = seat
	}

	if params.SendMethod != "" && params.SourceCode != "" {
		values[fmt.Sprintf("%s_send_method", params.SourceCode)] = params.SendMethod
	}

	if params.TrolleyToken != "" {
		values["trolley_token"] = params.TrolleyToken
	}

	if params.UserCommission {
		values["req_predicted_commission"] = "1"
	}

	for k, v := range params.Universal() {
		values[k] = v
	}

	return values
}
