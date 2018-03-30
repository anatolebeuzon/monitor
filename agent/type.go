package agent

import "time"

// Poller is the interface that wraps the Probe method.
type Poller interface {
	Poll() error
}

type Website struct {
	// Hostname string
	URL     string
	Metrics []Metric
}

type Websites []Website

type Metric struct {
	Date        time.Time
	DNStime     time.Duration
	TLStime     time.Duration
	ConnectTime time.Duration
	TTFB        time.Duration
	StatusCode  int
}
