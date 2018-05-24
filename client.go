package ticketswitch

import (
	"bytes"
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

func (client *Client) setHeaders(r *Request) error {
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
	return nil
}

// Do makes a request to the API
func (client *Client) Do(req *Request) (resp *http.Response, err error) {
	u, err := client.getURL(req)
	if err != nil {
		return
	}
	err = client.setHeaders(req)
	if err != nil {
		return
	}

	var body io.Reader
	if req.Body != nil {
		data, err := marshal(req.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(data)
		req.Header.Set("Content-Type", "application/json")
	}

	r, err := http.NewRequest(req.Method, u.String(), body)
	if err != nil {
		return
	}
	r.Header = req.Header

	resp, err = client.HTTPClient.Do(r)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = checkForError(resp)
	}
	return
}

// Test tests the API connection returning a User on success
func (client *Client) Test() (*User, error) {
	req := NewRequest("GET", "test.v1", nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

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
func (client *Client) ListEvents(params *ListEventsParams) (*ListEventsResults, error) {
	req := NewRequest(http.MethodGet, "events.v1", nil)
	if params != nil {
		req.SetValues(params.Params())
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

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
func (client *Client) GetEvents(eventIDs []string, params *UniversalParams) (map[string]*Event, error) {
	req := NewRequest(http.MethodGet, "events_by_id.v1", nil)
	if params != nil {
		req.SetValues(params.Universal())
	}

	req.SetValues(map[string]string{"event_id_list": strings.Join(eventIDs, ",")})

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

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
func (client *Client) GetEvent(eventID string, params *UniversalParams) (*Event, error) {
	events, err := client.GetEvents([]string{eventID}, params)

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
func (client *Client) ListPerformances(params *ListPerformancesParams) (*ListPerformancesResults, error) {
	req := NewRequest(http.MethodGet, "performances.v1", nil)
	if params != nil {
		req.SetValues(params.Params())
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var doc ListPerformancesTopLevel
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc.Results, nil
}

// GetAvailability fetches availability for a performce from the API
func (client *Client) GetAvailability(perf string, params *GetAvailabilityParams) (*AvailabilityResult, error) {
	req := NewRequest(http.MethodGet, "availability.v1", nil)
	if params != nil {
		req.SetValues(params.Params())
	}
	req.Values.Set("perf_id", perf)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var results AvailabilityResult
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

// Get the available sources
func (client *Client) GetSources(params *UniversalParams) (*SourcesResult, error) {
	req := NewRequest(http.MethodGet, "sources.v1", nil)
	if params != nil {
		req.SetValues(params.Universal())
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var sources []Source
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&sources)
	if err != nil {
		return nil, err
	}

	sourcesResult := &SourcesResult{Sources: sources}

	return sourcesResult, nil
}
