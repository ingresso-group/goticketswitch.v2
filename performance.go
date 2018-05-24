package ticketswitch

import "time"

// Performance describes and occurance of an Event.
type Performance struct {
	// identifier for the performance.
	ID string `json:"perf_id"`

	// the name of the performance.
	Name string `json:"perf_name"`

	// identifier for the event.
	EventID string `json:"event_id"`

	// the localised date and time for the performance.
	Datetime time.Time `json:"iso8601_date_and_time"`
	// a human readable description of the date of the performance.
	DateDesc string `json:"date_desc"`
	// a human readable description of the time of the performance.
	TimeDesc string `json:"time_desc"`
	// the number of minutes the performance is expected to run for.
	RunningTime int `json:"running_time"`

	// the performance has pool seats available.
	HasPoolSeats bool `json:"has_pool_seats"`
	// the performance has limited availability.
	IsLimited bool `json:"is_limited"`
	// the performance is a ghost performance and is nolonger available.
	IsGhost bool `json:"is_ghost"`

	// the maximum number of seats available to book in a single order. This
	// value is cached and may not be accurate.
	CachedMaxSeats int `json:"cached_max_seats"`

	// pricing summary, may also include offers.
	CostRange CostRange `json:"cost_range"`
	// pricing summary when no leaving single seats, may also include offers.
	NoSinglesCostRange CostRange `json:"no_singles_cost_range"`

	// summerised availability data for the performance. This data is cached
	// from previous availability calls and may not be accurate.
	AvailabilityDetails AvailabilityDetails `json:"avail_details"`
}

// ListPerformancesResults represents the results from a ListPerformance call
type ListPerformancesResults struct {
	// indicates that the related performances have names
	HasPerfNames bool `json:"has_perf_names"`

	// the current status of the pagination of the result set
	PagingStatus PagingStatus `json:"paging_status"`

	// performances returned by the call
	Performances []Performance `json:"performance"`
}

// ListPerformancesTopLevel is the top level of the json response from a ListPerformance call
type ListPerformancesTopLevel struct {
	// indicates that the performance list will contain only one performance
	// and this performance should be automatically selected for the customer.
	AutoSelect bool `json:"autoselect_this_performance"`

	// contains the ListPerformancesResults
	Results ListPerformancesResults `json:"results"`
}
