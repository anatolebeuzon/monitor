package client

import (
	"go-project-3/types"
	"time"
)

const pollingInterval = 1

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

func (scheduler *scheduler) init() *types.AggregateMapByURL {
	// Create receiver
	agg := make(types.AggregateMapByURL)

	go scheduler.receive(&agg)

	// Launch check routines
	for _, stat := range scheduler.config.Statistics {
		go func(stat Statistic) {
			scheduler.GetData(stat.Timespan)
			for range time.Tick(time.Duration(stat.Frequency) * time.Second) {
				scheduler.GetData(stat.Timespan)
			}
		}(stat)
	}

	return &agg
}

func (scheduler *scheduler) receive(agg *types.AggregateMapByURL) {
	for {
		datum := <-scheduler.data
		for _, item := range datum.Agg {
			_, ok := (*agg)[item.URL]
			if !ok {
				(*agg)[item.URL] = make(types.AggregateMapByTimespan)
			}
			(*agg)[item.URL][datum.Timespan] = item.Metrics
		}
		scheduler.updateUI <- true
	}
}
