package ticketswitch

import "github.com/shopspring/decimal"

type Seat struct {
	ColumnID         string `json:"column_id"`
	FullID           string `json:"full_id"`
	IsRestrictedView bool   `json:"is_restricted_view"`
	RowID            string `json:"row_id"`
}

type TicketOrder struct {
	DiscountCode       string          `json:"discount_code"`
	DiscountDesc       string          `json:"discount_desc"`
	NumberOfSeats      int             `json:"no_of_seats"`
	SaleSeatprice      decimal.Decimal `json:"sale_seatprice"`
	SaleSurcharge      decimal.Decimal `json:"sale_surcharge"`
	Seats              []Seat          `json:seats`
	TotalSaleSeatprice decimal.Decimal `json:"total_sale_seatprice"`
	TotalSaleSurcharge decimal.Decimal `json:"total_sale_surcharge"`
}

type TicketOrdersHolder struct {
	TicketOrders []TicketOrder `json:"ticket_order"`
}

type Order struct {
	Event                 Event              `json:"event"`
	GotRequestedSeats     bool               `json:"got_requested_seats"`
	ItemNumber            int                `json:"item_number"`
	Performance           Performance        `json:"performance"`
	PriceBandCode         string             `json:"price_band_code"`
	RequestedSeatIDs      []string           `json:"requested_seat_ids"`
	ReserveFailureComment string             `json:"reserve_failure_comment"`
	SeatRequestStatus     string             `json:"seat_request_status"`
	SendMethod            SendMethod         `json:"send_method"`
	TicketOrdersHolder    TicketOrdersHolder `json:"ticket_orders"`
	TicketTypeCode        string             `json:"ticket_type_code"`
	TicketTypeDesc        string             `json:"ticket_type_desc"`
	TotalNumberOfSeats    int                `json:"total_no_of_seats"`
	TotalSaleSeatprice    decimal.Decimal    `json:"total_sale_seatprice"`
	TotalSaleSurcharge    decimal.Decimal    `json:"total_sale_surcharge"`
}

type Bundle struct {
	OrderCount     int             `json:"bundle_order_count"`
	SourceCode     string          `json:"bundle_source_code"`
	SourceDesc     string          `json:"bundle_source_desc"`
	TotalCost      decimal.Decimal `json:"bundle_total_cost"`
	TotalSeatprice decimal.Decimal `json:"bundle_total_seatprice"`
	TotalSendCost  decimal.Decimal `json:"bundle_total_send_cost"`
	TotalSurcharge decimal.Decimal `json:"bundle_total_surcharge"`
	CurrencyCode   string          `json:"currency_code"`
	Orders         []Order         `json:"order"`
}

type Trolley struct {
	Bundles         []Bundle `json:"bundle"`
	TransactionUUID string   `json:"transaction_uuid"`
	BundleCount     int      `json:"trolley_bundle_count"`
	OrderCount      int      `json:"trolley_order_count"`
}
