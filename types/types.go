package types

import (
	"time"
)

type Package struct {
	Timespan int
	Websites []WebsiteMetric
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
