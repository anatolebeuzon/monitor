package payload

import (
	"time"
)

// Stats represents what is sent from the daemon to the client
type Stats struct {
	Timespan int
	Metrics  map[string]Metric
}

func NewStats(timespan int) Stats {
	return Stats{
		Timespan: timespan,
		Metrics:  make(map[string]Metric),
	}
}

type Availabilities []WebsiteAvailability

type WebsiteAvailability struct {
	URL          string
	Date         time.Time
	Availability float64
}

type WebsiteMetric struct {
	URL    string
	Metric Metric
}

type Metric struct {
	AvgAvail         float64
	MinTTFB          time.Duration
	MaxTTFB          time.Duration
	AvgTTFB          time.Duration
	StatusCodeCounts map[int]int
}

func (m Metric) String() string {
	return "Average TTFB: " + m.AvgTTFB.String() + "\nMin TTFB: " + m.MinTTFB.String() + "\nMax TTFB: " + m.MaxTTFB.String()
}
