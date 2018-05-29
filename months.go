package ticketswitch

type Month struct {
	Month                string `json:"month"`                  //: "feb",
	MonthDesc            string `json:"month_desc"`             //: "February",
	MonthDatesBitmask    int    `json:"month_dates_bitmask"`    //: 251526135,
	MonthWeekdaysBitmask int    `json:"month_weekdays_bitmask"` //: 63,
	Year                 int    `json:"year"`                   //: 2017
}

// AvailabilityResult describes the current state of available seats for a
// Performance
type MonthsResult struct {
	Months []Month `json:"month"`
}

// GetAvailabilityParams are parameters that can be passed to the
// GetAvailability call.
type GetMonthsParams struct {
	UniversalParams
	EventID         string
	CostRange       bool
	CostRangeDetail bool
	BestValueOffer  bool
	MaxSavingOffer  bool
	MinCostOffer    bool
	TopPriceOffer   bool
	NoSinglesData   bool
}

// Params returns the call parameters as a map
func (params *GetMonthsParams) Params() map[string]string {
	values := make(map[string]string)

	if params.CostRange {
		values["req_cost_range"] = "1"
	}

	if params.BestValueOffer {
		values["req_best_value_offer"] = "1"
	}

	if params.CostRangeDetail {
		values["req_cost_range_details"] = "1"
	}

	if params.MaxSavingOffer {
		values["req_cost_range_max_saving_offer"] = "1"
	}

	if params.MinCostOffer {
		values["req_cost_range_min_cost_offer"] = "1"
	}

	if params.TopPriceOffer {
		values["req_cost_range_top_price_offer"] = "1"
	}

	if params.NoSinglesData {
		values["req_cost_range_no_singles_data"] = "1"
	}

	for k, v := range params.Universal() {
		values[k] = v
	}

	return values
}
