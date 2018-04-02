package client

import (
	"monitor/payload"
	"time"
)

type Scheduler struct {
	Config   Config
	Received receivers
	UpdateUI chan bool
}

type receivers struct {
	stats  chan payload.Stats
	alerts chan payload.Alerts
}

func NewScheduler(c Config) *Scheduler {
	return &Scheduler{
		Config: c,
		Received: receivers{
			stats:  make(chan payload.Stats),
			alerts: make(chan payload.Alerts),
		},
		UpdateUI: make(chan bool),
	}
}

func (s *Scheduler) Init(store *Store) {
	// Create receiver
	go s.receive(store)

	// Launch stat check routines
	c := &s.Config
	for _, stat := range []Statistic{c.Statistics.Left, c.Statistics.Right} {
		go func(stat Statistic) {
			s.GetData(stat.Timespan)
			for range time.Tick(time.Duration(stat.Frequency) * time.Second) {
				s.GetData(stat.Timespan)
			}
		}(stat)
	}

	// Launch alert check routine
	go func() {
		for range time.Tick(time.Duration(c.Alerts.Frequency) * time.Second) {
			s.GetAlerts(c.Alerts.Timespan)
		}
	}()
}

func (s *Scheduler) receive(store *Store) {
	for {
		select {
		case stats := <-s.Received.stats:
			for url, metric := range stats.Metrics {
				// Check that URL is registered
				if _, ok := store.Metrics[url]; !ok {
					store.Metrics[url] = make(map[int]Metric)
					store.URLs = append(store.URLs, url)
				}

				history := store.Metrics[url][stats.Timespan].AvgRespHist
				start := 0
				itemsToKeep := 30
				if len(history) >= itemsToKeep {
					start = len(history) - itemsToKeep + 1
				}
				history = append(history[start:], metric.Average.Response)

				store.Metrics[url][stats.Timespan] = Metric{
					Latest:      metric,
					AvgRespHist: history,
				}
			}
			s.UpdateUI <- true

		case alerts := <-s.Received.alerts:
			for url, alert := range alerts {
				// TODO: no check that URL is registered. Is it okay?
				store.Alerts[url] = append(store.Alerts[url], alert)
			}
			s.UpdateUI <- true
		}
	}
}
