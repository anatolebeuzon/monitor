package client

import (
	"monitor/payload"
	"strconv"
)

// Store contains all the data needed by the dashboard
type Store struct {
	URLs    []string
	Metrics Metrics
	Alerts  Alerts
}

// Metrics[url][timespan] will give the aggregatedMetric for the selected URL and timespan
type Metrics map[string]map[int]payload.Metric

type Alerts map[string][]payload.Alert

func NewStore() *Store {
	return &Store{
		URLs:    []string{},
		Metrics: make(Metrics),
		Alerts:  make(Alerts),
	}
}

func (a Alerts) String(url string) (str string) {
	for _, alert := range a[url] {
		str += "Website " + url + " is "
		if alert.BelowThreshold {
			str += "down. "
		} else {
			str += "up. "
		}
		str += "availability=" + strconv.FormatFloat(alert.Availability, 'f', 3, 64)
		str += ", time=" + alert.Date.String() + "\n"
	}
	return
}
