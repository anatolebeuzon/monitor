/*
Not OK, fix comments about startIdx

This file contains the logic regarding poll results aggregation, such as:
- getting the poll results of the last n seconds
- computing average availability
- computing average and maximum times (response times, TLS handshake times, etc.)
- counting HTTP response codes and counting client errors
*/

package daemon

import (
	"monitor/payload"
	"time"
)

// Aggregate returns a payload.Metric containing the statistics for the website,
// aggregated over the specified timespan in seconds.
func (w *Website) Aggregate(tf payload.Timeframe) payload.Metric {
	// Copy poll results to ensure that they are not modified by
	// concurrent functions while results are being aggregated
	// TODO: avoid this somehow?
	p := w.PollResults.Extract(tf)
	return payload.Metric{
		Availability:     p.Availability(),
		Average:          p.Average(),
		Max:              p.Max(),
		StatusCodeCounts: p.CountCodes(),
		ErrorCounts:      p.CountErrors(),
	}
}

// StartIndexFor (timespan) returns the index (startIndex) of the first
// trace result that is included in the provided timespan (in seconds).
// In other words, t[startIndex:] will be the metrics obtained between [now, now - timespan].
//
// It leverages the fact that poll results are sorted by increasing date.
// The returned startIdx can then be used to aggregate the metrics fetched
// during the specified timespan.
//
// For example, given the following PollResults:
//		[]PollResult{
//			{ currentTime - 6 minutes, ... }
//			{ currentTime - 4 minutes, ... }
//			{ currentTime - 2 minutes, ... }
//			{ currentTime, ... },
//		}
// and given timespan = 180 (seconds), StartIndexFor(timespan) would return 2,
// as it is the index of the first PollResult of the slice
// that occured in the timeframe [now, now - 180 seconds]
func (p PollResults) Extract(tf payload.Timeframe) PollResults {
	// Traverse the slice from the end to the beginning
	// (generally faster, as p might be a very long slice
	// if a large number of poll results are retained)
	var startIdx, endIdx int
	for i := len(p) - 1; i >= 0; i-- {
		if endIdx == 0 && p[i].Date.Before(tf.EndDate) {
			// if endIdx hasn't been set yet and if the current date is in the timeframe
			endIdx = i + 1
		}
		if p[i].Date.Before(tf.StartDate) {
			startIdx = i + 1
			break
		}
	}
	r := p[startIdx:endIdx]
	return r
}

// Availability returns the average availability based on the latest poll results,
// starting from startIdx. The return value is between 0 and 1.
func (p PollResults) Availability() float64 {
	if len(p) == 0 {
		// No poll result is available in the timeframe, so
		// we cannot know whether the website is up or down.
		// In this case, act as if the website is down.
		return float64(0)
	}

	c := 0
	for _, r := range p {
		if r.IsValid() {
			c++
		}
	}
	return float64(c) / float64(len(p))
}

// IsValid returns whether the poll result is considered valid or not.
//
// To be considered valid, the associated request must satisfy two conditions:
// the request did not end with an error, and
// the HTTP response code is neither a Client error nor a Server error.
func (p PollResult) IsValid() bool {
	return (p.Error == nil) && (p.StatusCode < 400)
}

// Average returns a payload.Timing, in which each duration (DNS, TCP, TLS...)
// is the average of the respective durations of the selected poll results
// (here, "selected" poll results refer to the poll results from index startIdx onwards).
func (p PollResults) Average() (avg payload.Timing) {

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
// is the maximum of the respective durations of the selected poll results
// (here, "selected" poll results refer to the poll results from index startIdx onwards).
func (p PollResults) Max() (max payload.Timing) {
	for _, r := range p {
		max.DNS = maxDuration(r.Timing.DNS, max.DNS)
		max.TCP = maxDuration(r.Timing.TCP, max.TCP)
		max.TLS = maxDuration(r.Timing.TLS, max.TLS)
		max.Server = maxDuration(r.Timing.Server, max.Server)
		max.TTFB = maxDuration(r.Timing.TTFB, max.TTFB)
		max.Transfer = maxDuration(r.Timing.Transfer, max.Transfer)
		max.Response = maxDuration(r.Timing.Response, max.Response)
	}
	return
}

// maxDuration returns the maximum duration of two durations.
func maxDuration(d1, d2 time.Duration) time.Duration {
	if d1 > d2 {
		return d1
	}
	return d2
}

// CountCodes counts the HTTP response codes in the latest poll results, starting from startIdx.
// The return value maps from each HTTP response code encountered to the number of such codes.
func (p PollResults) CountCodes() map[int]int {
	codesCount := make(map[int]int)
	for _, r := range p {
		if r.StatusCode != 0 {
			// If the request led to an HTTP response code, and not an error
			codesCount[r.StatusCode]++
		}
	}
	return codesCount
}

// CountErrors counts the errors in the latest poll results, starting from startIdx.
// The return value maps from each error string encountered to the number of such errors.
func (p PollResults) CountErrors() map[string]int {
	errorsCount := make(map[string]int)
	for _, r := range p {
		if r.Error != nil {
			errorsCount[r.Error.Error()]++
		}
	}
	return errorsCount
}
