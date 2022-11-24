package ticketswitch

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// SortMostPopular sorts results based on sales across all partners over
	// the last 48 hours
	SortMostPopular = "most_popular"
	// SortAlphabetic sorts results alphabetically by event description
	SortAlphabetic = "alphabetic"
	// SortCostAscending sorts results by minimum total price with the lowest
	// price first
	SortCostAscending = "cost_ascending"
	// SortCostDescending sorts results by maximum total price with the highest
	// price first
	SortCostDescending = "cost_descending"
	// SortCriticRating sorts results by the average critic rating with the
	// highest rating first
	SortCriticRating = "critic_rating"
	// SortRecent sorts results by the date they were added to the system with
	// the newest result first.
	SortRecent = "recent"
	// SortLastSale sorts results by the products that sold the most recently.
	SortLastSale = "last_sale"
)

type FunctionParams interface {
	Params() map[string]string
}

// Client wraps the ticketswitch f13 API.
type Client struct {
	Config     *Config
	HTTPClient *http.Client
}

// NewClient returns a pointer to a newly created client.
func NewClient(config *Config) *Client {
	client := Client{
		Config:     config,
		HTTPClient: http.DefaultClient,
	}
	return &client
}

func (client *Client) getURL(r *Request) (*url.URL, error) {
	u, err := url.Parse(client.Config.BaseURL)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	for k, vs := range r.Values {
		for _, v := range vs {
			q.Add(k, v)
		}
	}

	if client.Config.CryptoBlock != "" {
		if client.Config.User == "" {
			return nil, fmt.Errorf("ticketswitch: config specifies cryptoblock but doesn't supply a user")
		}
		q.Set("user_id", client.Config.User)
		q.Set("crypto_block", client.Config.CryptoBlock)
	}

	if client.Config.SubUser != "" {
		q.Set("sub_id", client.Config.SubUser)
	}

	u.RawQuery = q.Encode()
	u.Path = fmt.Sprintf("%s/f13/%s", u.Path, r.Endpoint)
	return u, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (client *Client) setHeaders(ctx context.Context, r *Request) error {
	if client.Config.Language != "" {
		r.Header.Set("Accept-Language", client.Config.Language)
	}

	if client.Config.CryptoBlock == "" {
		if client.Config.User == "" {
			return fmt.Errorf("ticketswitch: config does not specify a user")
		}
		if client.Config.Password == "" {
			return fmt.Errorf("ticketswitch: config does not specify a password")
		}

		r.Header.Set("Authorization", "Basic "+basicAuth(client.Config.User, client.Config.Password))
	}

	// Set a session tracking id if provided in context
	trackingId, ok := GetSessionTrackingID(ctx)
	if ok {
		r.Header.Set("x-request-id", trackingId)
	}
	return nil
}

// Do make a request to the API
func (client *Client) Do(ctx context.Context, req *Request) (resp *http.Response, err error) {
	u, err := client.getURL(req)
	if err != nil {
		return nil, err
	}
	err = client.setHeaders(ctx, req)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if req.Body != nil {
		data, err2 := marshal(req.Body)
		if err2 != nil {
			return nil, err2
		}
		body = bytes.NewBuffer(data)
		req.Header.Set("Content-Type", "application/json")
	}

	r, err := http.NewRequest(req.Method, u.String(), body)
	if err != nil {
		return nil, err
	}
	r.Header = req.Header
	r = r.WithContext(ctx)

	resp, err = client.HTTPClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = checkForError(resp)
		return resp, err
	}

	return resp, nil
}

// Test tests the API connection returning a User on success
func (client *Client) Test(ctx context.Context) (*User, error) {
	req := NewRequest("GET", "test.v1", nil)
	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UniversalParams are parameters that can be passed to any call
type UniversalParams struct {
	Availability                 bool
	AvailabilityWithPerformances bool
	ExtraInfo                    bool
	Reviews                      bool
	Media                        bool
	CostRange                    bool
	BestValueOffer               bool
	MaxSavingOffer               bool
	MinCostOffer                 bool
	TopPriceOffer                bool
	NoSinglesData                bool
	CostRangeDetails             bool
	SourceInfo                   bool
	TrackingID                   string
	Misc                         map[string]string
}

// Universal returns the parameters as a map of parameters
func (params *UniversalParams) Universal() map[string]string {
	v := make(map[string]string)

	if params.Availability {
		v["req_avail_details"] = "1"
	}
	if params.AvailabilityWithPerformances {
		v["req_avail_details"] = "1"
		v["req_avail_details_with_perfs"] = "1"
	}
	if params.ExtraInfo {
		v["req_extra_info"] = "1"
	}
	if params.Reviews {
		v["req_reviews"] = "1"
	}
	if params.Media {
		v["req_media_triplet_one"] = "1"
		v["req_media_triplet_two"] = "1"
		v["req_media_triplet_three"] = "1"
		v["req_media_triplet_four"] = "1"
		v["req_media_triplet_five"] = "1"
		v["req_media_seating_plan"] = "1"
		v["req_media_square"] = "1"
		v["req_media_landscape"] = "1"
		v["req_media_marquee"] = "1"
		v["req_video_iframe"] = "1"
	}
	if params.CostRange {
		v["req_cost_range"] = "1"
	}
	if params.BestValueOffer {
		v["req_cost_range"] = "1"
		v["req_best_value_offer"] = "1"
	}
	if params.MaxSavingOffer {
		v["req_cost_range"] = "1"
		v["req_max_saving_offer"] = "1"
	}
	if params.MinCostOffer {
		v["req_cost_range"] = "1"
		v["req_min_cost_offer"] = "1"
	}
	if params.TopPriceOffer {
		v["req_cost_range"] = "1"
		v["req_top_price_offer"] = "1"
	}
	if params.NoSinglesData {
		v["req_cost_range"] = "1"
		v["req_no_singles_data"] = "1"
	}
	if params.CostRangeDetails {
		v["req_cost_range_details"] = "1"
	}
	if params.SourceInfo {
		v["req_src_info"] = "1"
	}
	if params.TrackingID != "" {
		v["custom_tracking_id"] = params.TrackingID
	}

	for k, val := range params.Misc {
		v[k] = val
	}
	return v
}

// PaginationParams are parameters that can be passed to any call that
// paginates it's response
type PaginationParams struct {
	PageLength int
	PageNumber int
}

// Pagination returns the pagination parameters as a map.
func (params *PaginationParams) Pagination() map[string]string {
	v := make(map[string]string)
	if params.PageLength > 0 {
		v["page_len"] = strconv.Itoa(params.PageLength)
	}
	if params.PageNumber > 0 {
		v["page_no"] = strconv.Itoa(params.PageNumber)
	}
	return v
}

// ListEventsParams are parameters that can be passed to the ListEvents call.
type ListEventsParams struct {
	UniversalParams
	PaginationParams
	Keywords    []string
	StartDate   time.Time
	EndDate     time.Time
	CountryCode string
	CityCode    string
	Circle      *Circle
	IncludeDead bool
	SortOrder   string
}

// Params returns the call parameters as a map
func (params *ListEventsParams) Params() map[string]string {
	values := make(map[string]string)

	if len(params.Keywords) > 0 {
		values["keywords"] = strings.Join(params.Keywords, ",")
	}

	if dr := DateRange(params.StartDate, params.EndDate); dr != "" {
		values["date_range"] = dr
	}

	if params.CountryCode != "" {
		values["country_code"] = params.CountryCode
	}

	if params.CityCode != "" {
		values["city_code"] = params.CityCode
	}

	if params.Circle != nil && params.Circle.Valid() {
		values["circle"] = params.Circle.Param()
	}

	if params.IncludeDead {
		values["include_dead"] = "1"
	}

	if params.SortOrder != "" {
		values["sort_order"] = params.SortOrder
	}

	for k, v := range params.Universal() {
		values[k] = v
	}

	for k, v := range params.Pagination() {
		values[k] = v
	}

	return values
}

// ListEvents returns a paginated slice of Events from the API.
func (client *Client) ListEvents(ctx context.Context, params *ListEventsParams) (*ListEventsResults, error) {
	req := NewRequest(http.MethodGet, "events.v1", nil)
	if params != nil {
		req.SetValues(params.Params())
	}

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var doc map[string]json.RawMessage
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&doc)
	if err != nil {
		return nil, err
	}

	rawResults, ok := doc["results"]
	if !ok {
		return nil, errors.New("ticketswitch: no results in ListEvents response")
	}

	var results ListEventsResults
	err = json.Unmarshal(rawResults, &results)

	if err != nil {
		return nil, err
	}

	rawCurrencies, ok := doc["currency_details"]

	if ok {
		var currencies map[string]Currency
		err := json.Unmarshal(rawCurrencies, &currencies)

		if err != nil {
			return nil, err
		}
		results.Currencies = currencies
	}

	return &results, nil
}

type wrappedEvent struct {
	Event *Event `json:"event"`
}

type getEventResults struct {
	EventsByID map[string]wrappedEvent `json:"events_by_id"`
}

// GetEvents returns a map of events index by event ID from the API.
func (client *Client) GetEvents(ctx context.Context, eventIDs []string, params *UniversalParams) (map[string]*Event, error) {
	req := NewRequest(http.MethodGet, "events_by_id.v1", nil)
	if params != nil {
		req.SetValues(params.Universal())
	}

	req.SetValues(map[string]string{"event_id_list": strings.Join(eventIDs, ",")})

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results getEventResults
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&results)
	if err != nil {
		return nil, err
	}

	events := make(map[string]*Event)
	for k, v := range results.EventsByID {
		events[k] = v.Event
	}

	return events, nil
}

// ErrEventNotFound will be returned when a specific event has been requested
// but didn't get a result back in response with that ID
var ErrEventNotFound = errors.New("ticketswitch: event not found")

// GetEvent returns an Event fetched from the API
func (client *Client) GetEvent(ctx context.Context, eventID string, params *UniversalParams) (*Event, error) {
	events, err := client.GetEvents(ctx, []string{eventID}, params)

	if err != nil {
		return nil, err
	}

	event, ok := events[eventID]
	if !ok {
		return nil, ErrEventNotFound
	}

	return event, nil
}

// ListPerformancesParams are parameters that can be passed to the
// ListPerformances call.
type ListPerformancesParams struct {
	UniversalParams
	PaginationParams
	EventID   string
	StartDate time.Time
	EndDate   time.Time
}

// Params returns the call parameters as a map
func (params *ListPerformancesParams) Params() map[string]string {
	values := make(map[string]string)

	if params.EventID != "" {
		values["event_id"] = params.EventID
	}

	if dr := DateRange(params.StartDate, params.EndDate); dr != "" {
		values["date_range"] = dr
	}

	for k, v := range params.Pagination() {
		values[k] = v
	}

	for k, v := range params.Universal() {
		values[k] = v
	}

	return values
}

// ListPerformances fetches a slice of performances from the API
//
//nolint:dupl
func (client *Client) ListPerformances(ctx context.Context, params *ListPerformancesParams) (*ListPerformancesResults, error) {
	req := NewRequest(http.MethodGet, "performances.v1", nil)
	if params != nil {
		req.SetValues(params.Params())
	}

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var doc ListPerformancesTopLevel
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc.Results, nil
}

// ListPerformanceTimes fetches a slice of unique performance times from the API
//
//nolint:dupl
func (client *Client) ListPerformanceTimes(ctx context.Context, params *ListPerformancesParams) (*ListPerformanceTimesResults, error) {
	req := NewRequest(http.MethodGet, "times.v1", nil)
	if params != nil {
		req.SetValues(params.Params())
	}

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var doc ListPerformanceTimesTopLevel
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&doc)
	if err != nil {
		return nil, err
	}
	return &doc.Results, nil
}

// GetAvailability fetches availability for a performce from the API
//
//nolint:dupl
func (client *Client) GetAvailability(ctx context.Context, perf string, params *GetAvailabilityParams) (*AvailabilityResult, error) {
	req := NewRequest(http.MethodGet, "availability.v1", nil)
	if params != nil {
		req.SetValues(params.Params())
	}
	req.Values.Set("perf_id", perf)

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results AvailabilityResult
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

// GetDiscounts fetches the Discounts for a particular performance, ticket type and price band from the API
func (client *Client) GetDiscounts(ctx context.Context, perf, ticketTypeCode, priceBandCode string, params *UniversalParams) (*DiscountsResult, error) {
	req := NewRequest(http.MethodGet, "discounts.v1", nil)
	if params != nil {
		req.SetValues(params.Universal())
	}
	req.Values.Set("perf_id", perf)
	req.Values.Set("ticket_type_code", ticketTypeCode)
	req.Values.Set("price_band_code", priceBandCode)

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results DiscountsResult
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

// GetSources fetches the available sources (a.k.a. backend systems) from the
// API
func (client *Client) GetSources(ctx context.Context, params *UniversalParams) (*SourcesResult, error) {
	req := NewRequest(http.MethodGet, "sources.v1", nil)
	if params != nil {
		req.SetValues(params.Universal())
	}

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sources []Source
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&sources)
	if err != nil {
		return nil, err
	}

	sourcesResult := &SourcesResult{Sources: sources}

	return sourcesResult, nil
}

// nolint:dupl
// GetSendMethods fetches the available send methods for a performance from the
// API
func (client *Client) GetSendMethods(ctx context.Context, perf string, params *UniversalParams) (*SendMethodsResults, error) {
	req := NewRequest(http.MethodGet, "send_methods.v1", nil)
	if params != nil {
		req.SetValues(params.Universal())
	}
	req.Values.Set("perf_id", perf)

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results SendMethodsResults
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

// MakeReservation places a hold on products in the inventory via the API
func (client *Client) MakeReservation(ctx context.Context, params *MakeReservationParams) (*ReservationResult, error) {
	req := NewRequest(http.MethodPost, "reserve.v1", params.Params())

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var reservation ReservationResult
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&reservation)
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

// TransactionParams are parameters that can be passed into the
// ReleaseReservation and GetStatus calls.
type TransactionParams struct {
	UniversalParams
	TransactionUUID string
}

// Params returns the call parameters as a map
func (params *TransactionParams) Params() map[string]string {
	values := map[string]string{
		"transaction_uuid": params.TransactionUUID,
	}

	for k, v := range params.Universal() {
		values[k] = v
	}
	return values
}

// ReleaseReservation makes a best effort attempt to release any reservations
// made on backend systems for a transaction.
func (client *Client) ReleaseReservation(ctx context.Context, params *TransactionParams) (success bool, err error) {
	req := NewRequest(http.MethodPost, "release.v1", params.Params())

	resp, err := client.Do(ctx, req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result map[string]bool
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return false, err
	}

	ok, success := result["released_ok"]

	if !ok {
		success = false
	}

	return success, nil
}

// MakePurchaseParams are the parameters that are passed into the MakePurchase
// call. A purchase must include the transaction UUID for an existing reserved
// transaction and some customer information. Optionally a payment method can
// be specified to provide payment details to the API when not purchasing on
// credit. If you require the API to send a confirmation email then set the
// SendConfirmationEmail flag to true (requires an email address to be
// specified in the customer information).
type MakePurchaseParams struct {
	UniversalParams
	TransactionUUID       string
	AgentReference        string
	Customer              Customer
	PaymentMethod         PaymentMethod
	SendConfirmationEmail bool
}

// Params returns the parameters needed to make the purchase call.
func (params *MakePurchaseParams) Params() map[string]string {
	values := map[string]string{
		"transaction_uuid": params.TransactionUUID,
	}

	if params.AgentReference != "" {
		values["agent_reference"] = params.AgentReference
	}

	if params.SendConfirmationEmail {
		values["send_confirmation_email"] = "1"
	}

	for k, v := range params.Customer.Params() {
		values[k] = v
	}

	if params.PaymentMethod != nil {
		for k, v := range params.PaymentMethod.PaymentParams() {
			values[k] = v
		}
	}
	for k, v := range params.Universal() {
		values[k] = v
	}
	return values
}

// MakePurchase attempts to purchase a previously reserved transaction via the
// API
func (client *Client) MakePurchase(ctx context.Context, params *MakePurchaseParams) (*MakePurchaseResult, error) {
	req := NewRequest(http.MethodPost, "purchase.v1", params.Params())

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result MakePurchaseResult
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetStatus retrieves the transaction from the API
func (client *Client) GetStatus(ctx context.Context, params *TransactionParams) (*StatusResult, error) {
	req := NewRequest(http.MethodGet, "status.v1", nil)
	if params != nil {
		req.SetValues(params.Params())
	}

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result StatusResult
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type CancellationParams struct {
	UniversalParams
	TransactionUUID string
	CancelItemsList CancelItemsList
}

type CancelItemsList []int

func (items CancelItemsList) String() string {
	stringArray := make([]string, len(items))
	for idx, val := range items {
		stringArray[idx] = fmt.Sprint(val)
	}
	return strings.Join(stringArray, ",")
}

func (params *CancellationParams) Params() map[string]string {
	values := map[string]string{
		"transaction_uuid": params.TransactionUUID,
	}
	if len(params.CancelItemsList) > 0 {
		values["cancel_items_list"] = params.CancelItemsList.String()
	}

	for k, v := range params.Universal() {
		values[k] = v
	}
	return values
}

// Cancel cancels transactions via the API
func (client *Client) Cancel(ctx context.Context, params *CancellationParams) (*CancellationResult, error) {
	req := NewRequest(http.MethodPost, "cancel.v1", nil)
	if params != nil {
		req.SetValues(params.Params())
	}
	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	var result CancellationResult
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// EmailCheck will check whether email passed meets
// RFC822 standard, returning error if not.
func (client *Client) EmailCheck(ctx context.Context, params *EmailCheckParams) error {
	if params == nil || params.EmailAddress == "" {
		return errors.New("No email was provided for verification")
	}
	req := NewRequest(http.MethodGet, "email_check.v1", nil)
	req.SetValues(params.Params())

	resp, err := client.Do(ctx, req)

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	return err
}
