package types

import "time"

type AggregateMetric struct {
	URL              string
	AvgAvail         float64
	MinTTFB          time.Duration
	MaxTTFB          time.Duration
	AvgTTFB          time.Duration
	StatusCodeCounts map[int]int
}

type AggregateMetrics []AggregateMetric

func (metric *AggregateMetric) String() string {
	return "Average TTFB: " + metric.AvgTTFB.String() + "\nMin TTFB: " + metric.MinTTFB.String() + "\nMax TTFB: " + metric.MaxTTFB.String()
}
