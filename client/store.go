/*
This file contains the types used to store the information available
for display on the dashboard.
*/

package client

import (
	"sync"
	"time"

	"github.com/oxlay/monitor/payload"
)

// Store contains all the data needed by the dashboard.
//
// It is regularly updated by goroutines that communicate with the daemon,
// and read by the dashboard, hence the mutex lock to avoid concurrent r/w.
type Store struct {
	sync.RWMutex
	URLs       []string
	CurrentIdx int // Index of the currently displayed website (website order is defined by Store.URLs)
	Metrics    Metrics
	Alerts     Alerts
}

// Metrics maps from each website and timespan to the corresponding Metric object.
//
// Metrics[url][timespan] will give the aggregated metric for the selected URL and timespan
type Metrics map[string]map[int]Metric

// Metric represents, for a given website and a given aggregation timespan,
// the corresponding statistics as needed by the dashboard.
type Metric struct {
	Latest payload.Metric

	// AvgRespHist stores the history of the average response times
	// It is used as a data source for the dashboard's graphs
	AvgRespHist []time.Duration
}

// Alerts maps from a URL to the alerts of the corresponding website.
type Alerts map[string][]payload.Alert

// NewStore creates a new Store and returns a pointer to it.
func NewStore() *Store {
	return &Store{
		URLs:    []string{},
		Metrics: make(Metrics),
		Alerts:  make(Alerts),
	}
}
