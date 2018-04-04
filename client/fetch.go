package client

import (
	"log"
	"monitor/payload"
	"net/rpc"
	"time"
)

type Fetcher struct {
	Config Config
	Store  *Store

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

// Init initiates regular polling of the daemon.
func (f *Fetcher) Init() {
	// Launch stat check routines
	c := &f.Config
	for _, t := range []TimeConf{c.Statistics.Left, c.Statistics.Right} {
		go func(t TimeConf) {
			f.GetStats(t.Timespan)
			for range time.Tick(time.Duration(t.Frequency) * time.Second) {
				f.GetStats(t.Timespan)
			}
		}(t)
	}

	// Launch alert check routine
	go func() {
		for range time.Tick(time.Duration(c.Alerts.Frequency) * time.Second) {
			f.GetAlerts(c.Alerts.Timespan)
		}
	}()
}

// GetStats gets the latest websites Stats from the daemon via RPC.
func (f *Fetcher) GetStats(timespan int) {
	tf := payload.NewTimeframe(timespan)
	var stats payload.Stats
	if err := f.CallRPC("Handler.Stats", &tf, &stats); err != nil {
		log.Fatal(err)
	}

	s := f.Store
	s.Lock()
	for url, metric := range stats.Metrics {
		// Check that URL is registered
		if _, ok := s.Metrics[url]; !ok {
			s.Metrics[url] = make(map[int]Metric)
			s.URLs = append(s.URLs, url)
		}

		history := s.Metrics[url][stats.Timeframe.Seconds].AvgRespHist
		start := 0
		itemsToKeep := 30
		if len(history) >= itemsToKeep {
			start = len(history) - itemsToKeep + 1
		}
		history = append(history[start:], metric.Average.Response)

		s.Metrics[url][stats.Timeframe.Seconds] = Metric{
			Latest:      metric,
			AvgRespHist: history,
		}
	}
	s.Unlock()

	f.UpdateUI <- true
}

// GetAlerts gets the latest websites Alerts from the daemon via RPC.
func (f *Fetcher) GetAlerts(timespan int) {
	tf := payload.NewTimeframe(timespan)
	var alerts payload.Alerts
	if err := f.CallRPC("Handler.Alerts", &tf, &alerts); err != nil {
		log.Fatal(err)
	}

	s := f.Store
	s.Lock()
	for url, alert := range alerts {
		s.Alerts[url] = append(s.Alerts[url], alert)
	}
	s.Unlock()

	f.UpdateUI <- true
}

// CallRPC connects to the daemon, calls the named function, waits for
// it to complete, then closes the connection.
func (f *Fetcher) CallRPC(method string, args interface{}, reply interface{}) error {
	client, err := rpc.DialHTTP("tcp", f.Config.Server)
	if err != nil {
		return err
	}

	if err = client.Call(method, args, reply); err != nil {
		return err
	}

	return client.Close()
}
