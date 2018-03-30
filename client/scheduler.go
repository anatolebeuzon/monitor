package client

import (
	"go-project-3/types"
	"time"
)

type scheduler struct {
	config   Config
	data     chan types.AggregateByTimespan
	updateUI chan bool
}

func newScheduler(config Config) scheduler {
	return scheduler{
		config:   config,
		data:     make(chan types.AggregateByTimespan),
		updateUI: make(chan bool),
	}
}

func (s *scheduler) init() *types.AggregateMapByURL {
	// Create receiver
	agg := types.NewAggregateMapByURL()

	go s.receive(&agg)

	// Launch check routines
	for _, stat := range s.config.Statistics {
		go func(stat Statistic) {
			// TODO: this is not optimal, fix this
			s.GetData(stat.Timespan)
			for range time.Tick(time.Duration(stat.Frequency) * time.Second) {
				s.GetData(stat.Timespan)
			}
		}(stat)
	}

	return &agg
}

func (s *scheduler) receive(agg *types.AggregateMapByURL) {
	for {
		datum := <-s.data

		// Check that timespan is registered
		if _, ok := (*agg).TimespansLookup[datum.Timespan]; !ok {
			(*agg).TimespansLookup[datum.Timespan] = true
			(*agg).TimespansOrder = append((*agg).TimespansOrder, datum.Timespan)
		}

		for _, item := range datum.Agg {
			// Check that URL is registered
			if _, ok := (*agg).Map[item.URL]; !ok {
				(*agg).Map[item.URL] = make(types.AggregateMapByTimespan)
				(*agg).URLs = append((*agg).URLs, item.URL)
			}
			(*agg).Map[item.URL][datum.Timespan] = item.Metrics
		}
		s.updateUI <- true
	}
}
