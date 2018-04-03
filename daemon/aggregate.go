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
	"monitor/payload"
	"time"
)

// Aggregate returns a payload.Metric containing the statistics for the website,
// aggregated over the specified timespan in seconds.
func (w *Website) Aggregate(timespan int) payload.Metric {
	// Copy poll results to ensure that they are not modified by
	// concurrent functions while results are being aggregated
	// TODO: avoid this somehow?
	p := w.PollResults
	startIdx := p.StartIndexFor(timespan)
	return payload.Metric{
		Availability:     p.Availability(startIdx),
		Average:          p.Average(startIdx),
		Max:              p.Max(startIdx),
		StatusCodeCounts: p.CountCodes(startIdx),
		ErrorCounts:      p.CountErrors(startIdx),
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
func (p PollResults) StartIndexFor(timespan int) int {
	threshold := time.Now().Add(-time.Duration(timespan) * time.Second)

	// Traverse the slice from the end to the beginning
	// (generally faster, as p might be a very long slice
	// if a large number of poll results are retained)
	for i := len(p) - 1; i >= 0; i-- {
		if p[i].Date.Before(threshold) {
			return i + 1
		}
	}
	return 0
}

// Availability returns the average availability based on the latest poll results,
// starting from startIdx. The return value is between 0 and 1.
func (p PollResults) Availability(startIdx int) float64 {
	if len(p)-startIdx == 0 {
		// No recent enough poll result is available, so
		// we cannot know whether the website is up or down.
		// In this case, act as if the website is down.
		return float64(0)
	}

	c := 0
	for i := startIdx; i < len(p); i++ {
		if p[i].IsValid() {
			c++
		}
	}
	return float64(c) / float64(len(p)-startIdx)
}

// IsValid returns whether the poll result is considered valid or not.
//
// To be considered valid, the associated request must satisfy two conditions:
// the request did not end with an error, and
// the HTTP response code is neither a Client error nor a Server error.
func (p *PollResult) IsValid() bool {
	return (p.Error == nil) && (p.StatusCode < 400)
}

// Average returns a payload.Timing, in which each duration (DNS, TCP, TLS...)
// is the average of the respective durations of the selected poll results
// (here, "selected" poll results refer to the poll results from index startIdx onwards).
func (p PollResults) Average(startIdx int) (avg payload.Timing) {

	// Perform an attribute-wise sum of durations
	for i := startIdx; i < len(p); i++ {
		curr := p[i].Timing
		avg.DNS += curr.DNS
		avg.TCP += curr.TCP
		avg.TLS += curr.TLS
		avg.Server += curr.Server
		avg.TTFB += curr.TTFB
		avg.Transfer += curr.Transfer
		avg.Response += curr.Response
	}

	// Divide by the number of elements to get the average
	if len(p)-startIdx != 0 {
		n := time.Duration(len(p) - startIdx)
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
func (p PollResults) Max(startIdx int) (max payload.Timing) {
	for i := startIdx; i < len(p); i++ {
		curr := p[i].Timing
		max.DNS = maxDuration(curr.DNS, max.DNS)
		max.TCP = maxDuration(curr.TCP, max.TCP)
		max.TLS = maxDuration(curr.TLS, max.TLS)
		max.Server = maxDuration(curr.Server, max.Server)
		max.TTFB = maxDuration(curr.TTFB, max.TTFB)
		max.Transfer = maxDuration(curr.Transfer, max.Transfer)
		max.Response = maxDuration(curr.Response, max.Response)
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
func (p PollResults) CountCodes(startIdx int) map[int]int {
	codesCount := make(map[int]int)
	for i := startIdx; i < len(p); i++ {
		code := p[i].StatusCode
		if code != 0 {
			// If the request led to an HTTP response code, and not an error
			codesCount[code]++
		}
	}
	return codesCount
}

// CountErrors counts the errors in the latest poll results, starting from startIdx.
// The return value maps from each error string encountered to the number of such errors.
func (p PollResults) CountErrors(startIdx int) map[string]int {
	errorsCount := make(map[string]int)
	for i := startIdx; i < len(p); i++ {
		error := p[i].Error
		if error != nil {
			errorsCount[error.Error()]++
		}
	}
	return errorsCount
}
