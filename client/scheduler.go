package client

import (
	"monitor/payload"
	"time"
)

type Scheduler struct {
	Config   Config
	Received Receivers

	// UpdateUI informs the dashboard that new data is available in the store,
	// and that it should re-render to display the latest information
	UpdateUI chan bool
}

// Receivers are channels that temporarily store payloads from daemons,
// until those payloads are processed by Receive()
type Receivers struct {
	stats  chan payload.Stats
	alerts chan payload.Alerts
}

// NewScheduler creates a new Scheduler with the provided Config.
func NewScheduler(c Config) *Scheduler {
	return &Scheduler{
		Config: c,
		Received: Receivers{
			stats:  make(chan payload.Stats),
			alerts: make(chan payload.Alerts),
		},
		UpdateUI: make(chan bool),
	}
}

// Init initiates regular polling of the daemon.
func (s *Scheduler) Init(store *Store) {
	// Create receiver
	go s.Receive(store)

	// Launch stat check routines
	c := &s.Config
	for _, t := range []TimeConf{c.Statistics.Left, c.Statistics.Right} {
		go func(t TimeConf) {
			s.GetStats(t.Timespan)
			for range time.Tick(time.Duration(t.Frequency) * time.Second) {
				s.GetStats(t.Timespan)
			}
		}(t)
	}

	// Launch alert check routine
	go func() {
		for range time.Tick(time.Duration(c.Alerts.Frequency) * time.Second) {
			s.GetAlerts(c.Alerts.Timespan)
		}
	}()
}

// Receive processes received payloads and add those to the Store.
func (s *Scheduler) Receive(store *Store) {
	for {
		select {
		case stats := <-s.Received.stats:
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

		case alerts := <-s.Received.alerts:
			store.Lock()
			for url, alert := range alerts {
				store.Alerts[url] = append(store.Alerts[url], alert)
			}
			store.Unlock()
			s.UpdateUI <- true
		}
	}
}
