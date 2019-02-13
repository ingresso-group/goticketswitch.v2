package ticketswitch

// CancellationResult contains the results of the cancel API call.
type CancellationResult struct {
	CancelledItemNumbers []int               `json:"cancelled_item_numbers"`
	MustAlsoCancel       []Order             `json:"must_also_cancel"`
	Trolley              Trolley             `json:"trolley_contents"`
	CurrencyDetails      map[string]Currency `json:"currency_details"`
}
