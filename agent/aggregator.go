package agent

import (
	"go-project-3/payload"
)

func (w *Websites) Availability(timespan int) payload.Alerts {
	p := payload.NewAlerts()
	for _, website := range *w {
		p.Availabilities[website.URL] = website.Availability(timespan)
	}
	return p
}

func (w *Website) Availability(timespan int) float64 {
	// TODO: remove duplicated code with aggregateResults /!\

	// Copy trace results to ensure that they are not modified by
	// concurrent functions while results are being aggregated
	tr := w.TraceResults
	startIdx := tr.startIndexFor(timespan)
	return tr.Availability(startIdx)
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
		Availability:     tr.Availability(startIdx),
		MinTTFB:          minDuration(TTFBs),
		MaxTTFB:          maxDuration(TTFBs),
		AvgTTFB:          avgDuration(TTFBs),
		StatusCodeCounts: tr.CountCodes(startIdx),
	}
}
