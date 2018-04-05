/*
OK

This file contains the logic regarding poll results aggregation, such as:
- getting the poll results of the last n seconds
- computing average availability
- computing average and maximum times (response times, TLS handshake times, etc.)
- counting HTTP response codes and counting client errors
*/

package daemon

import (
	"time"

	"github.com/oxlay/monitor/internal/payload"
)

// Aggregate returns a payload.Metric containing the statistics for the website,
// aggregated over the specified timeframe in seconds.
func (w *Website) Aggregate(tf payload.Timeframe) payload.Metric {
	p := w.PollResults.Extract(tf)
	return payload.Metric{
		Availability:     Availability(p),
		Average:          Average(p),
		Max:              Max(p),
		StatusCodeCounts: CountCodes(p),
		ErrorCounts:      CountErrors(p),
	}
}

// Extract returns the poll results that are included in the provided timeframe.
//
// The returned poll results can then be used to aggregate the metrics fetched
// during the specified timeframe.
func (p *PollResults) Extract(tf payload.Timeframe) []PollResult {
	// Traverse the slice from the end to the beginning
	// (generally faster, as p is sorted by increasing date,
	// and only the latest poll results are generally wanted)
	var startIdx, endIdx int
	p.RLock()
	defer p.RUnlock()
	for i := len(p.items) - 1; i >= 0; i-- {
		if endIdx == 0 && p.items[i].Date.Before(tf.EndDate) {
			// if endIdx hasn't been set yet and if the current date is in the timeframe
			endIdx = i + 1
		}
		if p.items[i].Date.Before(tf.StartDate) {
			startIdx = i + 1
			break
		}
	}
	return p.items[startIdx:endIdx]
}

// Availability returns the average availability based on the provided poll results.
// The return value is between 0 and 1.
func Availability(p []PollResult) float64 {
	if len(p) == 0 {
		// No poll result is available, so  we cannot know
		// whether the website is up or down.
		// In this case, act as if the website is down.
		return float64(0)
	}

	// Compute availability
	c := 0
	for _, r := range p {
		if IsValid(r) {
			c++
		}
	}
	return float64(c) / float64(len(p))
}

// IsValid returns whether the poll result is considered valid or not.
//
// To be considered valid, the associated request must satisfy two criteria:
// the request did not end with an error, and
// the HTTP response code is neither a Client error nor a Server error.
func IsValid(p PollResult) bool {
	return (p.Error == nil) && (p.StatusCode < 400)
}

// Average returns a payload.Timing, in which each duration (DNS, TCP, TLS...)
// is the average of the respective durations in the poll results.
func Average(p []PollResult) (avg payload.Timing) {

	// Perform an attribute-wise sum of durations
	for _, r := range p {
		avg.DNS += r.Timing.DNS
		avg.TCP += r.Timing.TCP
		avg.TLS += r.Timing.TLS
		avg.Server += r.Timing.Server
		avg.TTFB += r.Timing.TTFB
		avg.Transfer += r.Timing.Transfer
		avg.Response += r.Timing.Response
	}

	// Divide by the number of elements to get the average
	if len(p) != 0 {
		n := time.Duration(len(p))
		avg.DNS /= n
		avg.TCP /= n
		avg.TLS /= n
		avg.Server /= n
		avg.TTFB /= n
		avg.Transfer /= n
		avg.Response /= n
	}

	return
}

// Max returns a payload.Timing, in which each duration (DNS, TCP, TLS...)
// is the maximum of the respective durations of the poll results.
func Max(p []PollResult) (max payload.Timing) {
	for _, r := range p {
		max.DNS = MaxDuration(r.Timing.DNS, max.DNS)
		max.TCP = MaxDuration(r.Timing.TCP, max.TCP)
		max.TLS = MaxDuration(r.Timing.TLS, max.TLS)
		max.Server = MaxDuration(r.Timing.Server, max.Server)
		max.TTFB = MaxDuration(r.Timing.TTFB, max.TTFB)
		max.Transfer = MaxDuration(r.Timing.Transfer, max.Transfer)
		max.Response = MaxDuration(r.Timing.Response, max.Response)
	}
	return
}

// MaxDuration returns the maximum duration of two durations.
func MaxDuration(d1, d2 time.Duration) time.Duration {
	if d1 > d2 {
		return d1
	}
	return d2
}

// CountCodes counts the HTTP response codes in the poll results.
// The return value maps from each HTTP response code encountered to the number of such codes.
func CountCodes(p []PollResult) map[int]int {
	codesCount := make(map[int]int)
	for _, r := range p {
		if r.StatusCode != 0 {
			// If the request led to an HTTP response code, and not an error
			codesCount[r.StatusCode]++
		}
	}
	return codesCount
}

// CountErrors counts the client, non-HTTP errors in the poll results.
// The return value maps from each error string encountered to the number of such errors.
func CountErrors(p []PollResult) map[string]int {
	errorsCount := make(map[string]int)
	for _, r := range p {
		if r.Error != nil {
			errorsCount[r.Error.Error()]++
		}
	}
	return errorsCount
}
