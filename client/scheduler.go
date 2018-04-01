package client

import (
	"go-project-3/payload"
	"time"
)

type scheduler struct {
	config   Config
	received chan payload.Stats
	updateUI chan bool
}

func newScheduler(c Config) *scheduler {
	return &scheduler{
		config:   c,
		received: make(chan payload.Stats),
		updateUI: make(chan bool),
	}
}

func (s *scheduler) init() *Store {
	// Create receiver
	store := NewStore()

	go s.receive(store)

	// Launch check routines
	for _, stat := range s.config.Statistics {
		go func(stat Statistic) {
			s.GetData(stat.Timespan)
			for range time.Tick(time.Duration(stat.Frequency) * time.Second) {
				s.GetData(stat.Timespan)
			}
		}(stat)
	}

	return store
}

func (s *scheduler) receive(store *Store) {
	for {
		stats := <-s.received

		// Check that timespan is registered
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
		s.updateUI <- true
	}
}
