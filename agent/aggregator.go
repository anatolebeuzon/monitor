package agent

import (
	"fmt"
	"go-project-3/types"
	"time"
)

func (websites *Websites) aggregateMetrics() (aggMet types.AggregateMetrics) {
	for _, website := range *websites {
		aggMet = append(aggMet, website.aggregateMetrics())
	}
	return
}

func (website *Website) aggregateMetrics() (aggMet types.AggregateMetric) {
	aggMet.URL = website.URL
	TTFBs := website.TTFBs()
	aggMet.AvgTTFB = avgDuration(TTFBs)
	aggMet.MinTTFB = minDuration(TTFBs)
	aggMet.MaxTTFB = maxDuration(TTFBs)
	aggMet.StatusCodeCounts = website.countCodes()
	return
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
	fmt.Println(codesCount)
	return codesCount
}
