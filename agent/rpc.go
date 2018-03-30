package agent

import (
	"context"
	"go-project-3/types"
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
func (handler *Handler) Websites(args int, reply *Websites) error {
	*reply = *handler.websites
	return nil
}

func (handler *Handler) Metrics(args int, reply *types.AggregateMetrics) error {
	*reply = handler.websites.aggregateMetrics()
	return nil
}

// StopDaemon stops the daemon.
// It operates by sending a stop signal to a channel.
// The stop signal will trigger a shutdown of the HTTP server, but will wait
// for all connections to return to idle.
// In particular, a "Handler.StopDaemon" call from a client will receive a response
// before the server shuts down.
func (handler *Handler) StopDaemon(_, _ *struct{}) error {
	handler.done <- true
	return nil
}

// ServeRPC starts an RPC server, and exposes the methods
// of the handler type.
func (websites *Websites) ServeRPC(port int) {
	done := make(chan bool)
	rpcServer := rpc.NewServer()
	rpcServer.Register(&Handler{websites: websites, done: done})
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

	httpServer.ListenAndServe()
}
