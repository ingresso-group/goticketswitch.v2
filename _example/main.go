package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	ticketswitch "github.com/ingresso-group/goticketswitch.v2"
)

func main() {
	setupTracing()

	config := ticketswitch.NewConfig("demo", "demopass")
	client := ticketswitch.NewClient(config)

	// Add opencensus tracing to the http client Transport (aka RoundTripper interface)
	client.HTTPClient.Transport = &ochttp.Transport{}
	parentCtx, span := trace.StartSpan(context.Background(), "Main")
	defer span.End()
	params := ticketswitch.GetAvailabilityParams{
		NumberOfSeats: 2,
	}

	var innerSpan *trace.Span
	// Add a timeout to the context
	ctx, cancel := context.WithTimeout(parentCtx, 1*time.Second)
	defer cancel()
	// Trace each individual API calls using inner spans
	ctx, innerSpan = trace.StartSpan(ctx, "GetAvailability")
	results, err := client.GetAvailability(ctx, "7AB-6", &params)
	innerSpan.End()
	fmt.Println("\n\nAVAILABILITY:")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", results)

	ctx, cancel = context.WithTimeout(parentCtx, 1*time.Second)
	defer cancel()
	ctx, innerSpan = trace.StartSpan(ctx, "GetDiscounts")
	discountsResults, err := client.GetDiscounts(ctx, "6IF-C5O", "CIRCLE", "A/pool", nil)
	innerSpan.End()
	fmt.Println("\n\ndiscountsResults:")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", discountsResults)

	ctx, cancel = context.WithTimeout(parentCtx, 3*time.Second)
	defer cancel()
	ctx, innerSpan = trace.StartSpan(ctx, "GetSources")
	sources, err := client.GetSources(ctx, nil)
	innerSpan.End()
	fmt.Println("\n\nSOURCES:")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", sources)

	ctx, cancel = context.WithTimeout(parentCtx, 1*time.Second)
	defer cancel()
	ctx, innerSpan = trace.StartSpan(ctx, "GetSendMethods")
	sendMethods, err := client.GetSendMethods(ctx, "7AB-5", nil)
	innerSpan.End()
	fmt.Println("\n\nSEND METHODS:")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", sendMethods)

	ctx, cancel = context.WithTimeout(parentCtx, 1*time.Second)
	defer cancel()
	reserveParams := &ticketswitch.MakeReservationParams{
		PerformanceID:  "7AB-5",
		PriceBandCode:  "B/pool",
		TicketTypeCode: "CIRCLE",
		NumberOfSeats:  2,
		Seats:          []string{"A1", "A2"},
		SourceCode:     sendMethods.SourceCode,
		SendMethod:     sendMethods.SendMethodsHolder.SendMethods[1].Code,
	}
	ctx, innerSpan = trace.StartSpan(ctx, "MakeReservation")
	reservation, err := client.MakeReservation(ctx, reserveParams)
	innerSpan.End()
	fmt.Println("\n\nRESERVATION:")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", reservation)

	ctx, cancel = context.WithTimeout(parentCtx, 1*time.Second)
	defer cancel()
	transactionParams := ticketswitch.TransactionParams{TransactionUUID: reservation.Trolley.TransactionUUID}
	ctx, innerSpan = trace.StartSpan(ctx, "GetStatus")
	status, err := client.GetStatus(ctx, &transactionParams)
	innerSpan.End()
	fmt.Println("\n\nSTATUS:")
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("%+v", status)

	ctx, cancel = context.WithTimeout(parentCtx, 1*time.Second)
	defer cancel()
	ctx, innerSpan = trace.StartSpan(ctx, "ReleaseReservation")
	success, err := client.ReleaseReservation(ctx, &transactionParams)
	innerSpan.End()
	fmt.Println("\n\nRELEASE:")
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("%+v", success)
}

// setupTracing is an example of how to set up tracing using opencensus.
// There are of course many ways to do this.
// You can run jaeger locally using the docker run command found here:
// https://www.jaegertracing.io/docs/getting-started/
func setupTracing() {
	// Register stats and trace exporters to export the collected data.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Tracing exporter via Jaeger
	jaegerHost := "http://localhost:14268"
	jeagerExporter, err := jaeger.NewExporter(jaeger.Options{
		Endpoint:    jaegerHost,
		ServiceName: "ticketswitcher",
	})
	if err != nil {
		log.Error(err)
	}
	defer jeagerExporter.Flush()
	trace.RegisterExporter(jeagerExporter)

	// Report stats at every second.
	view.SetReportingPeriod(1 * time.Second)
}
