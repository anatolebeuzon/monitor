package client

import (
	"fmt"
	"go-project-3/types"
	"log"
	"net/rpc"
)

const (
	rpcProtocol = "tcp"
	rpcServer   = "127.0.0.1:1234"
)

func GetData(agg *types.AggregateMetrics, receivedData chan bool) {
	client, err := rpc.DialHTTP(rpcProtocol, rpcServer)
	if err != nil {
		log.Fatal("Failed to connect to the daemon:", err)
	}
	args := 4
	err = client.Call("Handler.Metrics", &args, agg)
	if err != nil {
		log.Fatal("RPC error:", err)
	}
	receivedData <- true
	err = client.Close()
	if err != nil {
		log.Fatal("RPC closing error:", err)
	}
	// fmt.Printf("Websites: %v", reply)
}

func StopDaemon() {
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
