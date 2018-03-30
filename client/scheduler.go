package client

import (
	"go-project-3/types"
	"time"
)

const pollingInterval = 1

func schedulePolls(agg *types.AggregateMetrics, receivedData chan bool) {
	GetData(agg, receivedData)
	for range time.Tick(pollingInterval * time.Second) {
		GetData(agg, receivedData)
	}
}
