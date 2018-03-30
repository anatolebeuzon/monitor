package agent

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"
)

const retainedMetrics = 10

// Poll makes a GET request to a website, and measures response times and response codes.
func (website *Website) Poll() {
	req, err := http.NewRequest("GET", website.URL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	var start, connect, dns, tlsHandshake time.Time

	var metric Metric

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			metric.DNStime = time.Since(dns)
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			metric.TLStime = time.Since(tlsHandshake)
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			metric.ConnectTime = time.Since(connect)
		},

		GotFirstResponseByte: func() {
			metric.TTFB = time.Since(start)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	metric.StatusCode = resp.StatusCode
	metric.Date = time.Now()

	// Only retain the last metrics
	itemsToDelete := 0
	if len(website.Metrics) >= retainedMetrics {
		itemsToDelete = len(website.Metrics) + 1 - retainedMetrics
	}
	website.Metrics = append(website.Metrics[itemsToDelete:], metric)

	// fmt.Println(website)
	// fmt.Println(website.aggregateMetrics())
}

// isValid returns true if an HTTP return code is considered valid
// (i.e. not an HTTP error code)
func isValid(code int) bool {
	validCodes := []int{200, 301, 302}
	for _, validCode := range validCodes {
		if code == validCode {
			return true
		}
	}
	return false
}
