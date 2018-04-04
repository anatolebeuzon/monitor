package client

import (
	"monitor/payload"
	"sync"
	"time"
)

// Store contains all the data needed by the dashboard
type Store struct {
	sync.RWMutex
	URLs       []string
	currentIdx int // Index of the currently displayed website (website order is defined by Store.URLs)
	Metrics    Metrics
	Alerts     Alerts
}

// Metrics[url][timespan] will give the aggregated metric for the selected URL and timespan
type Metrics map[string]map[int]Metric

type Metric struct {
	Latest payload.Metric
	// AvgRespHist stores the history of the average response times
	AvgRespHist []time.Duration
}

type Alerts map[string][]payload.Alert

func NewStore() *Store {
	return &Store{
		URLs:    []string{},
		Metrics: make(Metrics),
		Alerts:  make(Alerts),
	}
}
