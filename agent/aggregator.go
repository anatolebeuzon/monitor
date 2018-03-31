package agent

import (
	"go-project-3/types"
	"time"
)

func (w *Websites) aggregateMetrics(timespan int) (p types.Payload) {
	p.Timespan = timespan
	for _, website := range *w {
		p.Websites = append(p.Websites, website.aggregateMetrics())
	}
	return
}

func (w *Website) aggregateMetrics() types.WebsiteMetric {
	TTFBs := w.TTFBs()
	return types.WebsiteMetric{
		URL: w.URL,
		Metric: types.Metric{
			MinTTFB:          minDuration(TTFBs),
			MaxTTFB:          maxDuration(TTFBs),
			AvgTTFB:          avgDuration(TTFBs),
			StatusCodeCounts: w.countCodes(),
		},
	}
}

func (w *Website) TTFBs() (durations []time.Duration) {
	for _, res := range w.TraceResults {
		durations = append(durations, res.TTFB)
	}
	return
}

func (w *Website) countCodes() map[int]int {
	codesCount := make(map[int]int)
	for _, res := range w.TraceResults {
		codesCount[res.StatusCode]++
	}
	// fmt.Println(codesCount)
	return codesCount
}
