package client

import (
	"fmt"
	"log"
	"monitor/payload"
	"net/rpc"
)

const rpcProtocol = "tcp"

func (s *Scheduler) GetData(timespan int) {
	client, err := rpc.DialHTTP(rpcProtocol, s.Config.Server)
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

	s.Received.stats <- stats
}

func (s *Scheduler) GetAlerts(timespan int) {
	// TODO: fix duplicated code with GetData()

	client, err := rpc.DialHTTP(rpcProtocol, s.Config.Server)
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

	s.Received.alerts <- alerts
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
