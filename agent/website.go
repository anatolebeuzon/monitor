package agent

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"
)

type Website struct {
	// Hostname string
	URL          string
	TraceResults TraceResults
}

// Poll makes a GET request to a website, and measures response times and response codes.
func (w *Website) Poll(retainedResults int) {
	req, err := http.NewRequest("GET", w.URL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	var start, connect, dns, tlsHandshake time.Time

	var res TraceResult

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			res.DNStime = time.Since(dns)
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			res.TLStime = time.Since(tlsHandshake)
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			res.ConnectTime = time.Since(connect)
		},

		GotFirstResponseByte: func() {
			res.TTFB = time.Since(start)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	res.StatusCode = resp.StatusCode
	res.Date = time.Now()

	// Only retain the last trace results
	// TODO: improve this
	itemsToDelete := 0
	if len(w.TraceResults) >= retainedResults {
		itemsToDelete = len(w.TraceResults) + 1 - retainedResults
	}
	w.TraceResults = append(w.TraceResults[itemsToDelete:], res)

	// fmt.Println(w)
	// fmt.Println(w.aggregateMetrics())
}

func (w *Website) schedulePolls(p PollConfig) {
	for range time.Tick(time.Duration(p.Interval) * time.Second) {
		w.Poll(p.RetainedResults)
	}
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
