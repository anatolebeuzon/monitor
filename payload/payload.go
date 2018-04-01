package payload

import (
	"time"
)

// Stats represents what is sent from the daemon to the client
type Stats struct { // TODO: add Date?
	Timespan int
	Metrics  map[string]Metric
}

func NewStats(timespan int) Stats {
	return Stats{
		Timespan: timespan,
		Metrics:  make(map[string]Metric), // map from a URL to a Metric
	}
}

type Alerts struct {
	Date           time.Time
	Availabilities map[string]float64 // map from a URL to an availability between 0 and 1
}

func NewAlerts() Alerts {
	return Alerts{
		Date:           time.Now(),
		Availabilities: make(map[string]float64),
	}
}

type Metric struct {
	Availability     float64
	MinTTFB          time.Duration
	MaxTTFB          time.Duration
	AvgTTFB          time.Duration
	StatusCodeCounts map[int]int
}

func (m Metric) String() string {
	return "Average TTFB: " + m.AvgTTFB.String() + "\nMin TTFB: " + m.MinTTFB.String() + "\nMax TTFB: " + m.MaxTTFB.String()
}
