/*
This file handles all the interactions with an RPC client.
It exposes both the aggregated statistics and the alerts.
*/

package agent

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

// Stats puts the latest websites stats, aggregated over the specified timespan,
// as the reply value.
//
// Stats is meant to be used through an RPC call.
func (h *Handler) Stats(timespan int, reply *payload.Stats) error {
	*reply = h.Websites.Stats(timespan)
	return nil
}

// Alerts puts the latest websites alerts as the reply value.
// An alert is created only if the availability threshold is crossed
// based on the availability aggregated over the specified timespan.
//
// Alerts is meant to be used through an RPC call.
func (h *Handler) Alerts(timespan int, reply *payload.Alerts) error {
	*reply = h.Websites.Alerts(timespan, h.AlertThreshold)
	return nil
}

// ServeRPC starts an RPC server, and publishes the methods
// of the Handler type.
func ServeRPC(h *Handler, port int, interrupt chan os.Signal) {
	// TODO: a bit more documentation here
	rpcServer := rpc.NewServer()
	rpcServer.Register(h)
	rpcServer.HandleHTTP("/_goRPC_", "/debug/rpc") // args are the defaults used by HandleHTTP

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: rpcServer,
	}

	// Gracefully handle shutdown requests
	// TODO: remove this? is a graceful shutdown really useful? or use TCP instead of HTTP?
	go func() {
		<-interrupt
		httpServer.Shutdown(context.Background())
	}()

	fmt.Println("Listening for RPC requests on port", port)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
