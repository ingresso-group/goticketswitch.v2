package ticketswitch

// PagingStatus describes the current status of the pagination of a result set
type PagingStatus struct {
	PageLength       int `json:"page_length"`
	PageNumber       int `json:"page_number"`
	PagesRemaining   int `json:"pages_remaining"`
	ResultsRemaining int `json:"results_remaining"`
	TotalResults     int `json:"total_unpaged_results"`
}
