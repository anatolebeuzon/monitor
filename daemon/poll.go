/*
OK

This file contains the polling logic, namely:
- how and when websites are polled
- how metrics are collected throughout the lifecycle of an HTTP request
- how poll results are saved (to allow for later analysis and aggregation)
*/

package daemon

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/oxlay/monitor/payload"
)

// SchedulePolls schedules regular polls for the website. It never returns.
// The polling interval and the metrics retention policy are defined by the config file.
func (w *Website) SchedulePolls() {
	for range time.Tick(time.Duration(w.Interval) * time.Second) {
		w.Poll()
	}
}

// Poll makes a GET request to a website, measuring various times
// throughout the HTTP request, and reading the HTTP response code.
func (w *Website) Poll() {

	// Create request
	req, err := http.NewRequest("GET", w.URL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Record the exact times when the different parts of the request are reached
	var t [7]time.Time // t will store those times
	trace := &httptrace.ClientTrace{
		DNSStart:             func(_ httptrace.DNSStartInfo) { t[0] = time.Now() },
		DNSDone:              func(_ httptrace.DNSDoneInfo) { t[1] = time.Now() },
		ConnectStart:         func(_, _ string) { t[2] = time.Now() },
		ConnectDone:          func(_, _ string, _ error) { t[3] = time.Now() },
		GotConn:              func(_ httptrace.GotConnInfo) { t[4] = time.Now() },
		GotFirstResponseByte: func() { t[5] = time.Now() },
	}

	// Execute request and read response
	var p PollResult
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := NewTransport().RoundTrip(req)
	if err != nil {
		p.Error = err
	} else {
		p.StatusCode = resp.StatusCode
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		t[6] = time.Now() // records the fact that the body has been read (response is over)
	}

	// If an error occured, some times of the t slice may still be at 0.
	// However, they must be set to a sensible time to avoid getting
	// absurd results when converting those times to meaningful durations.
	if t[0].IsZero() {
		t[0] = time.Now()
	}
	for i := range t {
		if (i > 0) && t[i].IsZero() {
			t[i] = t[i-1]
		}
	}

	p.Date = t[0]

	// Convert the recorded times to meaningful durations
	p.Timing = payload.Timing{
		DNS:      t[1].Sub(t[0]),
		TCP:      t[3].Sub(t[2]),
		TLS:      t[4].Sub(t[3]),
		Server:   t[5].Sub(t[4]),
		Transfer: t[6].Sub(t[5]),
		TTFB:     t[5].Sub(t[0]),
		Response: t[6].Sub(t[0]),
	}

	// Save the poll result at the end of the website's poll results
	w.SaveResult(&p)
}

// NewTransport creates a new http.Transport.
//
// It purposefully has low timeouts, to allow for quick error detection and alerting.
// Keep-alive is also disabled, to ensure that processes such as DNS lookup
// and TLS handshakes are tested at each request.
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

// SaveResult saves a PollResult at the end of a websites' PollResults.
//
// If the number of poll results exceeds the user-defined retainedResults parameter,
// the oldest items are deleted.
// If retainedResults = 0, no metric is ever deleted.
func (w *Website) SaveResult(p *PollResult) {
	w.PollResults.Lock()
	defer w.PollResults.Unlock()

	i := 0
	if (w.RetainedResults != 0) && (len(w.PollResults.items) >= w.RetainedResults) {
		i = len(w.PollResults.items) + 1 - w.RetainedResults
	}
	w.PollResults.items = append(w.PollResults.items[i:], *p)
}
