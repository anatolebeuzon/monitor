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

// Handler contains all the websites and provides methods
// exposed over RPC.
type Handler struct {
	Websites       *Websites
	AlertThreshold float64
}

func (h *Handler) Metrics(timespan int, reply *payload.Stats) error {
	*reply = h.Websites.aggregateResults(timespan)
	return nil
}

func (h *Handler) Alerts(timespan int, reply *payload.Alerts) error {
	*reply = h.Websites.Alerts(timespan, h.AlertThreshold)
	return nil
}

// ServeRPC starts an RPC server, and exposes the methods
// of the handler type.
func ServeRPC(h *Handler, port int, interrupt chan os.Signal) {
	rpcServer := rpc.NewServer()
	rpcServer.Register(h)
	rpcServer.HandleHTTP("/_goRPC_", "/debug/rpc") // use defaults used by HandleHTTP

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: rpcServer,
	}

	// Gracefully handle shutdown requests
	// TODO: remove this? is a graceful shutdown really useful?
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
