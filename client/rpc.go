package client

import (
	"fmt"
	"go-project-3/payload"
	"log"
	"net/rpc"
)

const rpcProtocol = "tcp"

func (s *scheduler) GetData(timespan int) {

	client, err := rpc.DialHTTP(rpcProtocol, s.config.Server)
	if err != nil {
		log.Fatal("Failed to connect to the daemon:", err)
	}

	var stats payload.Stats
	err = client.Call("Handler.Metrics", &timespan, &stats)
	if err != nil {
		log.Fatal("RPC error:", err)
	}

	err = client.Close()
	if err != nil {
		log.Fatal("RPC closing error:", err)
	}

	s.received.stats <- stats
}

func (s *scheduler) GetAlerts(timespan int) {
	// TODO: fix duplicated code with GetData()

	client, err := rpc.DialHTTP(rpcProtocol, s.config.Server)
	if err != nil {
		log.Fatal("Failed to connect to the daemon:", err)
	}

	var alerts payload.Alerts
	err = client.Call("Handler.Alerts", &timespan, &alerts)
	if err != nil {
		log.Fatal("RPC error:", err)
	}

	err = client.Close()
	if err != nil {
		log.Fatal("RPC closing error:", err)
	}

	s.received.alerts <- alerts
}

func StopDaemon(rpcServer string) {
	client, err := rpc.DialHTTP(rpcProtocol, rpcServer)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	err = client.Call("Handler.StopDaemon", struct{}{}, nil)
	if err != nil {
		log.Fatal("Failed to stop daemon:", err)
	}
	fmt.Println("Daemon stopped.")
}
