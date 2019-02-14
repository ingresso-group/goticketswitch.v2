package ticketswitch

// CancellationResult contains the results of the cancel API call.
type CancellationResult struct {
	CancelledItemNumbers []int               `json:"cancelled_item_numbers"`
	MustAlsoCancel       []Order             `json:"must_also_cancel"`
	Trolley              Trolley             `json:"trolley_contents"`
	CurrencyDetails      map[string]Currency `json:"currency_details"`
}

// IsFullyCancelled checks the CancellationResult to see if the cancellation
// successfully cancelled all orders within the Trolley. If some orders are
// cancelled and others aren't then this will return false.
func (result *CancellationResult) IsFullyCancelled() bool {
	cancelled := true
	// Check if the trolley is empty
	if len(result.Trolley.Bundles) == 0 {
		return false
	}
	for _, bundle := range result.Trolley.Bundles {
		for _, order := range bundle.Orders {
			cancelled = cancelled && order.CancellationStatus == "cancelled"
		}
	}
	return cancelled
}
