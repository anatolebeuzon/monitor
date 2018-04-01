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

// Alerts maps from a URL to an Alert
type Alerts map[string]Alert

type Alert struct {
	Date           time.Time
	Availability   float64
	BelowThreshold bool // To know wether it is a new alert or a recovery
}

func NewAlert(availability float64, belowThreshold bool) Alert {
	return Alert{
		Date:           time.Now(),
		Availability:   availability,
		BelowThreshold: belowThreshold,
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
