package agent

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"
)

type Website struct {
	// Hostname string
	URL          string
	TraceResults TraceResults
	// DownAlertSent is true if at the last alert check from the front-end,
	// the aggregate availability was below the threshold. Keeping this information:
	// - avoids sending repetitive "website is down!" alerts
	// - enables the sending of one "website is up!" alert upon website recovery
	DownAlertSent bool
}

// Poll makes a GET request to a website, and measures response times and response codes.
func (w *Website) Poll(retainedResults int) {
	// TODO: refactor this function!
	req, err := http.NewRequest("GET", w.URL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	var start, connect, dns, tlsHandshake time.Time

	var tr TraceResult

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			dns = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			tr.DNSTime = time.Since(dns)
		},

		ConnectStart: func(network, addr string) {
			// TODO: do sth with addr?
			connect = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			// TODO: do sth with addr and err?
			tr.ConnectTime = time.Since(connect)
		},

		// TODO: use GotConn ?

		TLSHandshakeStart: func() {
			tlsHandshake = time.Now()
		},
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			// TODO: do sth with cs and err?
			tr.TLSTime = time.Since(tlsHandshake)
		},

		GotFirstResponseByte: func() {
			tr.TTFB = time.Since(start)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	transport := &http.Transport{
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 2 * time.Second,
			DualStack: true,
		}).DialContext,
		IdleConnTimeout:     2 * time.Second,
		TLSHandshakeTimeout: 2 * time.Second,
	}
	resp, err := transport.RoundTrip(req)
	tr.Date = time.Now()
	if err != nil {
		tr.Error = err
	} else {
		tr.StatusCode = resp.StatusCode
	}
	fmt.Println(tr)
	// Only retain the last trace results
	// TODO: improve this
	itemsToDelete := 0
	if len(w.TraceResults) >= retainedResults {
		itemsToDelete = len(w.TraceResults) + 1 - retainedResults
	}
	w.TraceResults = append(w.TraceResults[itemsToDelete:], tr)

	// fmt.Println(w)
	// fmt.Println(w.aggregateMetrics())
}

func (w *Website) schedulePolls(p PollConfig) {
	for range time.Tick(time.Duration(p.Interval) * time.Second) {
		w.Poll(p.RetainedResults)
	}
}
