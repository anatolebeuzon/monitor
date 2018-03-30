package client

import "time"

type Website struct {
	Hostname string
	URL      string
	Metrics  []Metric
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
