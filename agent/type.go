package agent

import "time"

// Poller is the interface that wraps the Probe method.
type Poller interface {
	Poll() error
}

type Website struct {
	// Hostname string
	URL          string
	TraceResults []TraceResult
}

type Websites []Website

type TraceResult struct {
	Date        time.Time
	DNStime     time.Duration
	TLStime     time.Duration
	ConnectTime time.Duration
	TTFB        time.Duration
	StatusCode  int
}
