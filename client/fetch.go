/*
This file is used to fetch the latest stats and alerts from the daemon, via RPC.

It contains:
- the scheduling logic (when data is fetched)
- the actual fetch logic (how data is fetched)
- the receiving logic (how the received data is saved to the Store)
*/

package client

import (
	"log"
	"monitor/payload"
	"net/rpc"
	"time"
)

// Fetcher is used by fetch methods to effectively craft requests, save responses
// to the store, and signal the dashboard that new data is available.
type Fetcher struct {
	Config Config // Config dictates how often requests to the daemon should be made
	Store  *Store // Store contains all the data available for display

	// UpdateUI informs the dashboard that new data is available in the store,
	// and that it should re-render to display the latest information
	UpdateUI chan bool
}

// NewFetcher creates a new Fetcher with the provided Config.
func NewFetcher(c Config, s *Store) *Fetcher {
	return &Fetcher{
		Config:   c,
		Store:    s,
		UpdateUI: make(chan bool),
	}
}

// Init initiates regular polling of stats and alerts from the daemon.
func (f *Fetcher) Init() {
	// Launch stat check routines
	c := &f.Config
	for _, t := range []TimeConf{c.Statistics.Left, c.Statistics.Right} {
		go func(t TimeConf) {
			f.GetStats(t.Timespan) // Get stats on dashboard startup, without waiting
			for range time.Tick(time.Duration(t.Frequency) * time.Second) {
				f.GetStats(t.Timespan)
			}
		}(t)
	}

	// Launch alert check routine
	go func() {
		f.GetAlerts(c.Alerts.Timespan) // Get alerts on dashboard startup, without waiting
		for range time.Tick(time.Duration(c.Alerts.Frequency) * time.Second) {
			f.GetAlerts(c.Alerts.Timespan)
		}
	}()
}

// GetStats gets the latest websites Stats from the daemon via RPC.
func (f *Fetcher) GetStats(timespan int) {
	// Craft and send request
	tf := payload.NewTimeframe(timespan)
	var stats payload.Stats
	if err := f.CallRPC("Handler.Stats", &tf, &stats); err != nil {
		log.Fatal(err)
	}

	// Save the resulting stats to the store
	s := f.Store
	s.Lock()
	for url, metric := range stats.Metrics {
		// Check that the URL is registered
		if _, ok := s.Metrics[url]; !ok {
			// If not, initialize the corresponding map,
			// and add the url to s.URLs to make the website
			// accessible on the dashboard.
			s.Metrics[url] = make(map[int]Metric)
			s.URLs = append(s.URLs, url)
		}

		// Add the received response time to the "average response time" graph data
		history := s.Metrics[url][stats.Timeframe.Seconds].AvgRespHist
		start := 0
		itemsToKeep := 30 // Number of points on the "average response time" graph
		if len(history) >= itemsToKeep {
			// Remove older data if necessary
			start = len(history) - itemsToKeep + 1
		}
		history = append(history[start:], metric.Average.Response)

		// Save the resulting Metric to the store
		s.Metrics[url][stats.Timeframe.Seconds] = Metric{
			Latest:      metric,
			AvgRespHist: history,
		}
	}
	s.Unlock()

	f.UpdateUI <- true // tell dashboard to rerender
}

// GetAlerts gets the latest websites Alerts from the daemon via RPC.
func (f *Fetcher) GetAlerts(timespan int) {
	// Craft and send request
	tf := payload.NewTimeframe(timespan)
	var alerts payload.Alerts
	if err := f.CallRPC("Handler.Alerts", &tf, &alerts); err != nil {
		log.Fatal(err)
	}

	// Save the resulting alerts to the store
	s := f.Store
	s.Lock()
	for url, alert := range alerts {
		s.Alerts[url] = append(s.Alerts[url], alert)
	}
	s.Unlock()

	f.UpdateUI <- true // tell dashboard to rerender
}

// CallRPC connects to the daemon, calls the named function, waits for
// it to complete, then closes the connection.
func (f *Fetcher) CallRPC(method string, args interface{}, reply interface{}) error {
	// Connect to the server
	client, err := rpc.DialHTTP("tcp", f.Config.Server)
	if err != nil {
		return err
	}

	// Make the remote procedure call
	if err = client.Call(method, args, reply); err != nil {
		return err
	}

	return client.Close()
}
