package ticketswitch

// GeoData contains the longitude and latitude of the Venue
type GeoData struct {
	// latitude of the venue.
	Latitude float64 `json:"latitude"`
	// longitude of the venue.
	Longitude float64 `json:"longitude"`
}

type UpsellList struct {
	EventIds []string `json:"event_id"`
}

// Event represents a product in the ticketswitch system.
type Event struct {
	// the identifier for the event.
	ID string `json:"event_id"`
	// status of the event.
	Status string `json:"event_status"`
	// human-readable name for the event.
	Description string `json:"event_desc"`
	// the backend system from which the event originates.
	Source string `json:"source_desc"`
	// the internal code for the backend system.
	SourceCode string `json:"source_code"`
	// the type of the event.
	EventType string `json:"event_type"`
	// a human-readable description of the venue.
	Venue string `json:"venue_desc"`
	// a dictionary of class descriptions that the event belongs to keyed on
	// class identifier.
	Classes map[string]string `json:"classes"`
	// a list of filters that the event belongs to.
	Filters []string `json:"custom_filter"`
	// venue post code.
	Postcode string `json:"postcode"`
	// venue geographical data.
	GeoData GeoData `json:"geo_data"`
	// human-readable venue city.
	City string `json:"city_desc"`
	// venue city code
	CityCode string `json:"city_code"`
	// human-readable country name.
	Country string `json:"country_desc"`
	// ISO 3166-1 country code.
	CountryCode string `json:"country_code"`
	// maximum running time of a performance in minutes.
	MaxRunningTime int `json:"max_running_time"`
	// minimum running time of a performance in minutes.
	MinRunningTime int `json:"min_running_time"`
	// indicates that the performance time for this event is relevant and
	// should be shown.
	ShowPerformanceTime bool `json:"show_perf_time"`
	// indicates that the event has no performances.
	HasNoPerformances bool `json:"has_no_perfs"`
	// indicates the event is seated.
	IsSeated bool `json:"is_seated"`
	// indicates that ticket purchases for this event will require a departure
	// date.
	NeedsDepartureDate bool `json:"needs_departure_date"`
	// indicates that ticket purchases for this event will require a duration.
	NeedsDuration bool `json:"needs_duration"`
	// indicates that ticket purchases for this event will require a
	// performance id.
	NeedsPerformance bool `json:"needs_performance"`
	// list of related event id's for upselling.
	UpsellList UpsellList `json:"event_upsell_list"`
	// pricing summary from cached availability. Only present when requested.
	CostRange CostRange `json:"cost_range"`
	// pricing summary from cached availability. Only present when requested.
	NoSinglesCostRange CostRange `json:"no_singles_cost_range"`
	// summary pricing information broken down by availability. This is cached
	// data. Only present when requested.
	CostRangeDetails CostRangeDetails `json:"cost_range_details"`
	// indexed on content name. Only present when requested.
	Content map[string]Content `json:"content"`
	// fields indexed on field name. Only present when requested.
	Fields map[string]Field `json:"fields"`
	// event info in plain text. Only present when requested.
	EventInfo string `json:"event_info"`
	// event info as HTML. Only present when requested.
	EventInfoHTML string `json:"event_info_html"`
	// venue address in plain text. Only present when requested.
	VenueAddr string `json:"venue_addr"`
	// venue address as HTML. Only present when requested.
	VenueAddrHTML string `json:"venue_addr_html"`
	// venue info in plain text. Only present when requested.
	VenueInfo string `json:"venue_info"`
	// venue info as HTML. Only present when requested.
	VenueInfoHTML string `json:"venue_info_html"`
	// media items indexed on media name. Only present when requested.
	Media map[string]Media `json:"media"`
	// reviews of this product Only present when requested.
	Reviews []Review `json:"reviews"`
	// summary of critic review star rating.  rated from 1 (lowest) to 5
	// (highest).
	CriticReviewPercent float64 `json:"critic_review_percent"`
	// summary of availability details from cached data. Only
	// present when requested.
	AvailabilityDetails AvailabilityDetails `json:"availability_details"`
	// list of Event objects that comprise a meta event.
	ComponentEvents []Event `json:"component_events"`
	// list of valid qualities available for purchase. from cached data, only
	// available when requested by **get_events** or **get_event**.
	ValidQuantities []int `json:"valid_quantities"`
	// indicates that the event is an addon product to another event.
	IsAddon bool `json:"is_add_on"`
	// AreaCode is for internal use only
	AreaCode string `json:"area_code"`
	// Code is for internal use only
	Code string `json:"event_code"`
	// VenueCode is for internal use only
	VenueCode string `json:"venue_code"`
}

// ListEventsResults represents a set of events returned by the API
type ListEventsResults struct {
	// map of Currency objects
	Currencies map[string]Currency
	// when an object doesn't explicitly state a currency code, this code
	// should be assumed.
	DefaultCurrencyCode string
	// the code of currency the user is expecting
	DesiredCurrencyCode string

	// the current status of the pagination of the result set
	PagingStatus PagingStatus `json:"paging_status"`

	// performances returned by the call
	Events []Event `json:"event"`
}
