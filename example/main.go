package main

import (
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
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", results)

	sources, err := client.GetSources(nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", sources)
}
