package payload

import (
	"time"
)

// Stats contains, for a given timespan (in seconds), the aggregated
// poll results for all the websites polled by the daemon.
type Stats struct {
	Timeframe Timeframe         // Time window use to aggregate results
	Metrics   map[string]Metric // Maps from a website URL to a Metric
}

// A Metric contains the aggregated poll results of one website.
type Metric struct {
	Availability     float64        // Average availability
	Average          Timing         // Average HTTP lifecycle times
	Max              Timing         // Max HTTP lifecycle times
	StatusCodeCounts map[int]int    // Maps from an HTTP response code to the number of times it was encountered
	ErrorCounts      map[string]int // Maps from a client error string to the number of times it was encountered
}

// A Timing contains the durations of each phase of an HTTP request.
type Timing struct {
	// DNS is the duration of the DNS lookup.
	DNS time.Duration

	// TCP is the TCP connection time.
	TCP time.Duration

	// TLS is the duration of the TLS handshake, if applicable.
	// If the website was contacted over HTTP, TLS will be set to time.Duration(0).
	TLS time.Duration

	// Server is the server processing time. It measures the time it takes the server to
	// deliver the first response byte since a connection was established.
	Server time.Duration

	// Transfer is the transfer time of the response.
	// It starts when the first byte is received and ends when the last byte is received.
	Transfer time.Duration

	// TTFB is the time to first byte.
	// It is equal to the sum of all the previous durations except the Transfer time:
	// TTFB = DNS + TCP + TLS + Server
	TTFB time.Duration

	// Response is the response time.
	// It is equal to the sum of the TTFB and the Transfer time:
	// Response = DNS + TCP + TLS + Server + Transfer
	Response time.Duration
}
