package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	ticketswitch "github.com/ingresso-group/goticketswitch.v2"
)

func main() {
	config := ticketswitch.NewConfig("demo", "demopass")
	client := ticketswitch.NewClient(config)
	params := ticketswitch.GetAvailabilityParams{
		NumberOfSeats: 2,
	}
	//params.CostRange = true
	results, err := client.GetAvailability("7AA-5", &params)
	fmt.Println("\n\nAVAILABILITY:")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", results)

	sources, err := client.GetSources(nil)
	fmt.Println("\n\nSOURCES:")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", sources)

	reserveParams := &ticketswitch.MakeReservationParams{
		PerformanceID:  "7AB-5",
		PriceBandCode:  "B/pool",
		TicketTypeCode: "CIRCLE",
		NumberOfSeats:  2,
		Seats:          []string{"A1", "A2"},
	}
	reservation, err := client.MakeReservation(reserveParams)
	fmt.Println("\n\nRESERVATION:")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", reservation)

	transactionParams := ticketswitch.TransactionParams{TransactionUUID: reservation.Trolley.TransactionUUID}
	status, err := client.GetStatus(&transactionParams)
	fmt.Println("\n\nSTATUS:")
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("%+v", status)

	success, err := client.ReleaseReservation(&transactionParams)
	fmt.Println("\n\nRELEASE:")
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("%+v", success)
}
