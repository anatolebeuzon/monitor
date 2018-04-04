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
	tf := payload.NewTimeframe(timespan)
	var stats payload.Stats
	if err := s.CallRPC("Handler.Stats", &tf, &stats); err != nil {
		log.Fatal(err)
	}

	store := s.Store
	store.Lock()
	for url, metric := range stats.Metrics {
		// Check that URL is registered
		if _, ok := store.Metrics[url]; !ok {
			store.Metrics[url] = make(map[int]Metric)
			store.URLs = append(store.URLs, url)
		}

		history := store.Metrics[url][stats.Timeframe.Seconds].AvgRespHist
		start := 0
		itemsToKeep := 30
		if len(history) >= itemsToKeep {
			start = len(history) - itemsToKeep + 1
		}
		history = append(history[start:], metric.Average.Response)

		store.Metrics[url][stats.Timeframe.Seconds] = Metric{
			Latest:      metric,
			AvgRespHist: history,
		}
	}
	store.Unlock()

	s.UpdateUI <- true
}

// GetAlerts gets the latest websites Alerts from the daemon via RPC.
func (s *Scheduler) GetAlerts(timespan int) {
	tf := payload.NewTimeframe(timespan)
	var alerts payload.Alerts
	if err := s.CallRPC("Handler.Alerts", &tf, &alerts); err != nil {
		log.Fatal(err)
	}

	store := s.Store
	store.Lock()
	for url, alert := range alerts {
		store.Alerts[url] = append(store.Alerts[url], alert)
	}
	store.Unlock()

	s.UpdateUI <- true
}
