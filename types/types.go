package types

import (
	"strconv"
	"time"
)

type AggregateByTimespan struct {
	Timespan int
	Agg      []AggregateItem
}

type AggregateItem struct {
	URL     string
	Metrics AggregatedMetric
}

// AggregateMetric maps from an aggregation timespan to the corresponding
type AggregateMapByURL map[string]AggregateMapByTimespan

type AggregateMapByTimespan map[int]AggregatedMetric

type AggregatedMetric struct {
	AvgAvail         float64
	MinTTFB          time.Duration
	MaxTTFB          time.Duration
	AvgTTFB          time.Duration
	StatusCodeCounts map[int]int
}

func (metric AggregatedMetric) String() string {
	return "Average TTFB: " + metric.AvgTTFB.String() + "\nMin TTFB: " + metric.MinTTFB.String() + "\nMax TTFB: " + metric.MaxTTFB.String()
}

func (metrics AggregateMapByTimespan) String() (str string) {
	for timespan, metric := range metrics {
		str += "Aggregate over " + strconv.Itoa(timespan) + " seconds :\n"
		str += metric.String()
		str += "\n\n"
	}
	return
}
