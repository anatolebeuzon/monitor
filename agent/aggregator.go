package agent

import (
	"go-project-3/types"
)

func (w *Websites) aggregateResults(timespan int) (p types.Payload) {
	p.Timespan = timespan
	for _, website := range *w {
		p.Websites = append(p.Websites, website.aggregateResults(timespan))
	}
	return
}

func (w *Website) aggregateResults(timespan int) types.WebsiteMetric {
	// Copy trace results to ensure that they are not modified by
	// concurrent functions while results are being aggregated
	tr := w.TraceResults
	startIdx := tr.startIndexFor(timespan)
	TTFBs := tr.TTFBs(startIdx)
	return types.WebsiteMetric{
		URL: w.URL,
		Metric: types.Metric{
			MinTTFB:          minDuration(TTFBs),
			MaxTTFB:          maxDuration(TTFBs),
			AvgTTFB:          avgDuration(TTFBs),
			StatusCodeCounts: tr.CountCodes(startIdx),
		},
	}
}
