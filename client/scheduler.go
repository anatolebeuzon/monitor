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

func (s *Scheduler) Init() *Store {
	// Create receiver
	store := NewStore()

	go s.receive(store)

	// Launch stat check routines
	for _, stat := range s.Config.Statistics {
		go func(stat Statistic) {
			s.GetData(stat.Timespan)
			for range time.Tick(time.Duration(stat.Frequency) * time.Second) {
				s.GetData(stat.Timespan)
			}
		}(stat)
	}

	// Launch alert check routine
	go func() {
		for range time.Tick(time.Duration(s.Config.Alerts.Frequency) * time.Second) {
			s.GetAlerts(s.Config.Alerts.Timespan)
		}
	}()

	return store
}

func (s *Scheduler) receive(store *Store) {
	for {
		select {
		case stats := <-s.Received.stats:

			// Check that timespan is registered
			// TODO: do this while reading config, it makes no sense to do it here
			if _, ok := store.Timespans.Lookup[stats.Timespan]; !ok {
				store.Timespans.Lookup[stats.Timespan] = true
				store.Timespans.Order = append(store.Timespans.Order, stats.Timespan)
			}

			for url, metric := range stats.Metrics {
				// Check that URL is registered
				if _, ok := store.Metrics[url]; !ok {
					store.Metrics[url] = make(map[int]payload.Metric)
					store.URLs = append(store.URLs, url)
				}
				store.Metrics[url][stats.Timespan] = metric
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
