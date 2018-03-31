package agent

import (
	"context"
	"go-project-3/types"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
)

// Handler contains all the websites and provides methods
// exposed over RPC.
type Handler struct {
	websites *Websites
	done     chan bool
}

// Websites replies with all the websites (including metrics info).
func (h *Handler) Websites(args int, reply *Websites) error {
	*reply = *h.websites
	return nil
}

func (h *Handler) Metrics(timespan int, reply *types.Package) error {
	*reply = h.websites.aggregateMetrics(timespan)
	return nil
}

// StopDaemon stops the daemon.
// It operates by sending a stop signal to a channel.
// The stop signal will trigger a shutdown of the HTTP server, but will wait
// for all connections to return to idle.
// In particular, a "Handler.StopDaemon" call from a client will receive a response
// before the server shuts down.
func (h *Handler) StopDaemon(_, _ *struct{}) error {
	h.done <- true
	return nil
}

// ServeRPC starts an RPC server, and exposes the methods
// of the handler type.
func (w *Websites) ServeRPC(port int) {
	done := make(chan bool)
	rpcServer := rpc.NewServer()
	rpcServer.Register(&Handler{websites: w, done: done})
	rpcServer.HandleHTTP("/_goRPC_", "/debug/rpc") // use defaults used by HandleHTTP

	// TODO: fix duplicate server info with client
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: rpcServer,
	}

	// Gracefully handle shutdown requests
	go func() {
		<-done
		httpServer.Shutdown(context.Background())
	}()

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
