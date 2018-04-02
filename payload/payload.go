package payload

import (
	"time"
)

// Stats represents what is sent from the daemon to the client
type Stats struct { // TODO: add Date?
	Timespan int
	Metrics  map[string]Metric
}

func NewStats(timespan int) Stats {
	return Stats{
		Timespan: timespan,
		Metrics:  make(map[string]Metric), // map from a URL to a Metric
	}
}

// Alerts maps from a URL to an Alert
type Alerts map[string]Alert

type Alert struct {
	Date           time.Time
	Availability   float64
	BelowThreshold bool // To know wether it is a new alert or a recovery
}

func NewAlert(availability float64, belowThreshold bool) Alert {
	return Alert{
		Date:           time.Now(),
		Availability:   availability,
		BelowThreshold: belowThreshold,
	}
}

type Metric struct {
	Availability     float64
	Average          Timing
	Max              Timing
	StatusCodeCounts map[int]int
	ErrorCounts      map[string]int
}

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

	// TTFB is the time to first byte.
	// It is equal to the sum of all the previous durations.
	TTFB time.Duration

	// Transfer is the transfer time of the response.
	// It starts when the first byte is received and ends when the last byte is received.
	Transfer time.Duration

	// Response is the response time.
	// It is equal to the sum of the TTFB and the Transfer time.
	Response time.Duration
}

func (t *Timing) ToSlice() []time.Duration {
	return []time.Duration{t.DNS, t.TCP, t.TLS, t.Server, t.TTFB, t.Transfer, t.Response}
}
