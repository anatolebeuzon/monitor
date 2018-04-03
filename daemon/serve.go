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
	"monitor/payload"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

// Handler contains all the necessary data to satisfy RPC calls.
// It will be registered as the RPC receiver and its methods will be published.
type Handler struct {
	Websites       *Websites
	AlertThreshold float64
}

// Stats puts the latest websites stats (aggregated over
// the specified timespan in seconds) as the reply value.
//
// Stats is meant to be used through an RPC call.
func (h *Handler) Stats(timespan int, p *payload.Stats) error {
	*p = payload.NewStats(timespan)
	for _, website := range *h.Websites {
		(*p).Metrics[website.URL] = website.Aggregate(timespan)
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
func (h *Handler) Alerts(timespan int, a *payload.Alerts) error {
	*a = make(payload.Alerts)
	for i, website := range *h.Websites {
		// Get average availability
		startIdx := website.PollResults.StartIndexFor(timespan)
		avail := website.PollResults.Availability(startIdx)

		if (avail < h.AlertThreshold) && !website.DownAlertSent {
			// if the website is considered down but no alert for this event was sent yet
			// create a "website is down" alert
			(*a)[website.URL] = payload.NewAlert(avail, true)
			(*h.Websites)[i].DownAlertSent = true
		} else if (avail >= h.AlertThreshold) && website.DownAlertSent {
			// if the website is considered up but website was last reported down
			// create a "website has recovered" alert
			(*a)[website.URL] = payload.NewAlert(avail, false)
			(*h.Websites)[i].DownAlertSent = false
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
