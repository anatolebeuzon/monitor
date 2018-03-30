package agent

import (
	"go-project-3/types"
	"time"
)

func (websites *Websites) aggregateMetrics(timespan int) (aggMet types.AggregateByTimespan) {
	aggMet.Timespan = timespan
	for _, website := range *websites {
		aggMet.Agg = append(aggMet.Agg, website.aggregateMetrics())
	}
	return
}

func (website *Website) aggregateMetrics() types.AggregateItem {
	TTFBs := website.TTFBs()
	aggMet := types.AggregateItem{
		URL: website.URL,
		Metrics: types.AggregatedMetric{
			MinTTFB:          minDuration(TTFBs),
			MaxTTFB:          maxDuration(TTFBs),
			AvgTTFB:          avgDuration(TTFBs),
			StatusCodeCounts: website.countCodes(),
		},
	}
	return aggMet
}

func (website *Website) TTFBs() (durations []time.Duration) {
	for _, metric := range website.Metrics {
		durations = append(durations, metric.TTFB)
	}
	return
}

func (website *Website) countCodes() map[int]int {
	codesCount := make(map[int]int)
	for _, metric := range website.Metrics {
		codesCount[metric.StatusCode]++
	}
	// fmt.Println(codesCount)
	return codesCount
}
