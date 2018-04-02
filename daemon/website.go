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
	req, err := http.NewRequest("GET", w.URL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	var tr TraceResult

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), NewTrace(&tr)))
	resp, err := NewTransport().RoundTrip(req)
	tr.Date = time.Now()
	if err != nil {
		tr.Error = err
	} else {
		tr.StatusCode = resp.StatusCode
		resp.Body.Close()
	}

	w.SaveResult(&tr, retainedResults)
}

// SaveResult saves a TraceResult at the end of a websites' TraceResult.
//
// If the number of TraceResults exceeds the user-defined retainedResults parameter,
// the oldest items are deleted.
// If retainedResults = 0, no metric is ever deleted.
func (w *Website) SaveResult(tr *TraceResult, retainedResults int) {
	itemsToDelete := 0
	if (retainedResults != 0) && (len(w.TraceResults) >= retainedResults) {
		itemsToDelete = len(w.TraceResults) + 1 - retainedResults
	} // If retainedResults is set to 0, store an unlimited number of results

	w.TraceResults = append(w.TraceResults[itemsToDelete:], *tr)
}

// schedulePolls schedules regular polls for the website.
// The polling interval and the metrics retention policy are defined by the config file.
func (w *Website) schedulePolls(p PollConfig) {
	for range time.Tick(time.Duration(p.Interval) * time.Second) {
		w.Poll(p.RetainedResults)
	}
}

func NewTrace(tr *TraceResult) *httptrace.ClientTrace {
	var start, connect, dns, tlsHandshake time.Time
	start = time.Now()

	return &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { tr.DNSTime = time.Since(dns) },
		ConnectStart: func(network, addr string) {
			// TODO: do sth with addr?
			connect = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			// TODO: do sth with addr and err?
			tr.ConnectTime = time.Since(connect)
		},

		// TODO: use GotConn ?

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			// TODO: do sth with cs and err?
			tr.TLSTime = time.Since(tlsHandshake)
		},

		GotFirstResponseByte: func() { tr.TTFB = time.Since(start) },
	}
}

// NewTransport creates a new http.Transport.
//
// It purposefully has low timeouts, to allow for quick error detection and alerting.
// Keep-alive is also disabled, to ensure that processes such as DNS lookup
// and TLS handshakes are tested at each requested.
func NewTransport() *http.Transport {
	return &http.Transport{
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 2 * time.Second,
			DualStack: true,
		}).DialContext,
		IdleConnTimeout:     2 * time.Second,
		TLSHandshakeTimeout: 2 * time.Second,
	}
}
