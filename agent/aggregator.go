package agent

import (
	"go-project-3/payload"
)

func (w *Websites) Availability(timespan int) (p payload.Stats) {
	return
}

func (w *Website) Availability(timespan int) (p payload.Stats) {
	return
}

func (w *Websites) aggregateResults(timespan int) payload.Stats {
	p := payload.NewStats(timespan)
	for _, website := range *w {
		p.Metrics[website.URL] = website.aggregateResults(timespan)
	}
	return p
}

func (w *Website) aggregateResults(timespan int) payload.Metric {
	// Copy trace results to ensure that they are not modified by
	// concurrent functions while results are being aggregated
	tr := w.TraceResults
	startIdx := tr.startIndexFor(timespan)
	TTFBs := tr.TTFBs(startIdx)
	return payload.Metric{
		MinTTFB:          minDuration(TTFBs),
		MaxTTFB:          maxDuration(TTFBs),
		AvgTTFB:          avgDuration(TTFBs),
		StatusCodeCounts: tr.CountCodes(startIdx),
	}
}
