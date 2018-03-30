package client

import (
	"fmt"
	"go-project-3/types"
	"log"
	"net/rpc"
)

const rpcProtocol = "tcp"

func (s *scheduler) GetData(timespan int) {

	client, err := rpc.DialHTTP(rpcProtocol, s.config.Server)
	if err != nil {
		log.Fatal("Failed to connect to the daemon:", err)
	}

	var agg types.AggregateByTimespan
	err = client.Call("Handler.Metrics", &timespan, &agg)
	if err != nil {
		log.Fatal("RPC error:", err)
	}

	err = client.Close()
	if err != nil {
		log.Fatal("RPC closing error:", err)
	}

	s.data <- agg
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
