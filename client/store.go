package client

import (
	"go-project-3/payload"
	"strconv"
)

// Store stores metrics data
type Store struct {
	Timespans Timespans
	URLs      []string
	// Metrics[url][timespan] will give the aggregatedMetric for the selected URL and timespan
	Metrics map[string]map[int]payload.Metric
}

type Timespans struct {
	Order  []int
	Lookup map[int]bool
}

func NewStore() *Store {
	return &Store{
		Timespans: Timespans{
			Order:  []int{},
			Lookup: make(map[int]bool),
		},
		URLs:    []string{},
		Metrics: make(map[string]map[int]payload.Metric),
	}
}

func (s Store) String(url string) (str string) {
	for _, timespan := range s.Timespans.Order {
		str += "Aggregate over " + strconv.Itoa(timespan) + " seconds :\n"
		str += s.Metrics[url][timespan].String()
		str += "\n\n"
	}
	return
}
