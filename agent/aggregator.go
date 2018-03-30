package agent

import (
	"go-project-3/types"
	"time"
)

func (w *Websites) aggregateMetrics(timespan int) (aggMet types.AggregateByTimespan) {
	aggMet.Timespan = timespan
	for _, website := range *w {
		aggMet.Agg = append(aggMet.Agg, website.aggregateMetrics())
	}
	return
}

func (w *Website) aggregateMetrics() types.AggregateItem {
	TTFBs := w.TTFBs()
	aggMet := types.AggregateItem{
		URL: w.URL,
		Metrics: types.AggregatedMetric{
			MinTTFB:          minDuration(TTFBs),
			MaxTTFB:          maxDuration(TTFBs),
			AvgTTFB:          avgDuration(TTFBs),
			StatusCodeCounts: w.countCodes(),
		},
	}
	return aggMet
}

func (w *Website) TTFBs() (durations []time.Duration) {
	for _, metric := range w.Metrics {
		durations = append(durations, metric.TTFB)
	}
	return
}

func (w *Website) countCodes() map[int]int {
	codesCount := make(map[int]int)
	for _, metric := range w.Metrics {
		codesCount[metric.StatusCode]++
	}
	// fmt.Println(codesCount)
	return codesCount
}
