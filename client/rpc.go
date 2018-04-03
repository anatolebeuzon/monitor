package client

import (
	"log"
	"monitor/payload"
	"net/rpc"
)

// CallRPC connects to the daemon, calls the named function, waits for
// it to complete, then closes the connection.
func (s *Scheduler) CallRPC(method string, args interface{}, reply interface{}) error {
	client, err := rpc.DialHTTP("tcp", s.Config.Server)
	if err != nil {
		return err
	}

	if err = client.Call(method, args, reply); err != nil {
		return err
	}

	return client.Close()
}

// GetStats gets the latest websites Stats from the daemon via RPC.
func (s *Scheduler) GetStats(timespan int) {
	var stats payload.Stats
	if err := s.CallRPC("Handler.Stats", &timespan, &stats); err != nil {
		log.Fatal(err)
	}
	s.Received.stats <- stats
}

// GetAlerts gets the latest websites Alerts from the daemon via RPC.
func (s *Scheduler) GetAlerts(timespan int) {
	var alerts payload.Alerts
	if err := s.CallRPC("Handler.Alerts", &timespan, &alerts); err != nil {
		log.Fatal(err)
	}
	s.Received.alerts <- alerts
}
