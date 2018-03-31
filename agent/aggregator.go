package agent

import (
	"go-project-3/types"
)

func (w *Websites) aggregateMetrics(timespan int) (p types.Payload) {
	p.Timespan = timespan
	for _, website := range *w {
		p.Websites = append(p.Websites, website.aggregateMetrics(timespan))
	}
	return
}

func (w *Website) aggregateMetrics(timespan int) types.WebsiteMetric {
	startIdx := w.TraceResults.startIndexFor(timespan)
	TTFBs := w.TraceResults.TTFBs(startIdx)
	return types.WebsiteMetric{
		URL: w.URL,
		Metric: types.Metric{
			MinTTFB:          minDuration(TTFBs),
			MaxTTFB:          maxDuration(TTFBs),
			AvgTTFB:          avgDuration(TTFBs),
			StatusCodeCounts: w.TraceResults.CountCodes(startIdx),
		},
	}
}
