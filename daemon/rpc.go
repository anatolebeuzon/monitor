package agent

import (
	"context"
	"fmt"
	"log"
	"monitor/payload"
	"net/http"
	"net/rpc"
	"strconv"
)

// Handler contains all the websites and provides methods
// exposed over RPC.
type Handler struct {
	Websites       *Websites
	AlertThreshold float64
	Done           chan bool
}

func (h *Handler) Metrics(timespan int, reply *payload.Stats) error {
	*reply = h.Websites.aggregateResults(timespan)
	return nil
}

func (h *Handler) Alerts(timespan int, reply *payload.Alerts) error {
	*reply = h.Websites.Alerts(timespan, h.AlertThreshold)
	return nil
}

// StopDaemon stops the daemon.
// It operates by sending a stop signal to a channel.
// The stop signal will trigger a shutdown of the HTTP server, but will wait
// for all connections to return to idle.
// In particular, a "Handler.StopDaemon" call from a client will receive a response
// before the server shuts down.
func (h *Handler) StopDaemon(_, _ *struct{}) error {
	h.Done <- true
	return nil
}

// ServeRPC starts an RPC server, and exposes the methods
// of the handler type.
func ServeRPC(h *Handler, port int) {
	rpcServer := rpc.NewServer()
	rpcServer.Register(h)
	rpcServer.HandleHTTP("/_goRPC_", "/debug/rpc") // use defaults used by HandleHTTP

	// TODO: fix duplicate server info with client
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: rpcServer,
	}

	// Gracefully handle shutdown requests
	go func() {
		<-h.Done
		httpServer.Shutdown(context.Background())
	}()

	fmt.Println("Listening for RPC requests on port", port)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
