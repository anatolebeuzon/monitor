package client

import (
	"go-project-3/types"
	"time"
)

type scheduler struct {
	config   Config
	received chan types.Payload
	updateUI chan bool
}

func newScheduler(c Config) *scheduler {
	return &scheduler{
		config:   c,
		received: make(chan types.Payload),
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
		Package := <-s.received

		// Check that timespan is registered
		if _, ok := store.Timespans.Lookup[Package.Timespan]; !ok {
			store.Timespans.Lookup[Package.Timespan] = true
			store.Timespans.Order = append(store.Timespans.Order, Package.Timespan)
		}

		for _, w := range Package.Websites {
			// Check that URL is registered
			if _, ok := store.Metrics[w.URL]; !ok {
				store.Metrics[w.URL] = make(map[int]types.Metric)
				store.URLs = append(store.URLs, w.URL)
			}
			store.Metrics[w.URL][Package.Timespan] = w.Metric
		}
		s.updateUI <- true
	}
}
