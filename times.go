package ticketswitch

import (
	"time"
)

// PerformanceTime describes a possible date and time of a Performance.
type PerformanceTime struct {
	// the localized date and time for the performance.
	Datetime time.Time `json:"iso8601_date_and_time,omitempty"`

	// Human-readable description of performance time.
	TimeDesc string `json:"time_desc,omitempty"`
}

// ListPerformanceTimesResults represents the results from a ListPerformanceTimes call.
type ListPerformanceTimesResults struct {
	// Performance times returned by the call
	Times []PerformanceTime `json:"time,omitempty"`
}

// ListPerformanceTimesTopLevel represents the top level of the ListPerfanceTimes call's json response.
type ListPerformanceTimesTopLevel struct {
	Results ListPerformanceTimesResults `json:"results,omitempty"`
}
