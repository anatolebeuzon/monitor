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
type AggregateMapByURL struct {
	TimespansOrder  []int
	TimespansLookup map[int]bool
	URLs            []string
	Map             map[string]AggregateMapByTimespan
}

func NewAggregateMapByURL() AggregateMapByURL {
	return AggregateMapByURL{
		TimespansOrder:  []int{},
		TimespansLookup: make(map[int]bool),
		URLs:            []string{},
		Map:             make(map[string]AggregateMapByTimespan),
	}
}

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

func (agg AggregateMapByURL) String(url string) (str string) {
	for _, timespan := range agg.TimespansOrder {
		str += "Aggregate over " + strconv.Itoa(timespan) + " seconds :\n"
		str += agg.Map[url][timespan].String()
		str += "\n\n"
	}
	return
}
