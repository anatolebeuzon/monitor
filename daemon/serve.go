/*
OK

This file handles all the interactions with an RPC client.
It exposes both the aggregated statistics and the alerts.
*/

package daemon

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"strconv"

	"github.com/oxlay/monitor/payload"
)

// Handler contains all the necessary data to satisfy RPC calls.
// It will be registered as the RPC receiver and its methods will be published.
type Handler Websites

// Stats puts the latest websites stats (aggregated over
// the specified timespan in seconds) as the reply value.
//
// Stats is meant to be used through an RPC call.
func (h *Handler) Stats(tf payload.Timeframe, p *payload.Stats) error {
	*p = payload.Stats{Timeframe: tf, Metrics: make(map[string]payload.Metric)}
	for _, website := range *h {
		(*p).Metrics[website.URL] = website.Aggregate(tf)
	}
	return nil
}

// Alerts puts the latest websites alerts as the reply value.
//
// For each website, it compares its availability
// (on average, over the specified timespan in seconds)
// against the threshold, and creates an alert if the threshold is crossed.
//
// Alerts is meant to be used through an RPC call.
func (h *Handler) Alerts(tf payload.Timeframe, a *payload.Alerts) error {
	*a = make(payload.Alerts)
	for i, website := range *h {
		// Get average availability
		avail := Availability(website.PollResults.Extract(tf))

		if (avail < website.Threshold) && !website.DownAlertSent {
			// if the website is considered down but no alert for this event was sent yet
			// create a "website is down" alert
			(*a)[website.URL] = payload.Alert{Timeframe: tf, Availability: avail, BelowThreshold: true}
			(*h)[i].DownAlertSent = true
		} else if (avail >= website.Threshold) && website.DownAlertSent {
			// if the website is considered up but website was last reported down
			// create a "website has recovered" alert
			(*a)[website.URL] = payload.Alert{Timeframe: tf, Availability: avail, BelowThreshold: false}
			(*h)[i].DownAlertSent = false
		}
	}
	return nil
}

// ServeRPC starts an RPC server, and publishes the methods
// of the Handler type.
func ServeRPC(h *Handler, port int, interrupt chan os.Signal) {
	// Create RPC server
	rpcServer := rpc.NewServer()
	rpcServer.Register(h)                          // Publish Handler's methods
	rpcServer.HandleHTTP("/_goRPC_", "/debug/rpc") // args are the defaults used by HandleHTTP

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: rpcServer,
	}

	// Gracefully handle shutdown requests
	go func() {
		<-interrupt
		httpServer.Shutdown(context.Background())
	}()

	// Begin serving HTTP requests
	fmt.Println("Listening for RPC requests on port", port)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
