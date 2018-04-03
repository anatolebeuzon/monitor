package agent

import (
	"fmt"
	"io/ioutil"
	"monitor/payload"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"
)

type Website struct {
	URL         string
	PollResults PollResults

	// DownAlertSent is true if at the last alert check from the front-end,
	// the aggregate availability was below the threshold. Keeping this information:
	// - avoids sending repetitive "website is down!" alerts
	// - enables the sending of one "website is up!" alert upon website recovery
	DownAlertSent bool
}

// Poll makes a GET request to a website, measures response times and response codes.
func (w *Website) Poll(retainedResults int) {
	req, err := http.NewRequest("GET", w.URL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	var t [7]time.Time
	trace := &httptrace.ClientTrace{
		DNSStart:             func(_ httptrace.DNSStartInfo) { t[0] = time.Now() },
		DNSDone:              func(_ httptrace.DNSDoneInfo) { t[1] = time.Now() },
		ConnectStart:         func(_, _ string) { t[2] = time.Now() },
		ConnectDone:          func(_, _ string, err error) { t[3] = time.Now() }, // TODO: handle err?
		GotConn:              func(_ httptrace.GotConnInfo) { t[4] = time.Now() },
		GotFirstResponseByte: func() { t[5] = time.Now() },
	}

	var p PollResult
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := NewTransport().RoundTrip(req)
	if err != nil {
		p.Error = err
	} else {
		p.StatusCode = resp.StatusCode
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		t[6] = time.Now() // body has been read, response is over
	}

	// If an error occured, some times may still be at 0
	// However, they must be set to a sensible time to avoid getting absurd results
	if t[0].IsZero() {
		t[0] = time.Now()
	}
	for i := range t {
		if (i > 0) && t[i].IsZero() {
			t[i] = t[i-1]
		}
	}

	p.Date = t[0]
	p.Timing = payload.Timing{
		DNS:      t[1].Sub(t[0]),
		TCP:      t[3].Sub(t[2]),
		TLS:      t[4].Sub(t[3]),
		Server:   t[5].Sub(t[4]),
		Transfer: t[6].Sub(t[5]),
		TTFB:     t[5].Sub(t[0]),
		Response: t[6].Sub(t[0]),
	}

	// fmt.Println(p)
	w.SaveResult(&p, retainedResults)
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

// SaveResult saves a PollResult at the end of a websites' PollResults.
//
// If the number of poll results exceeds the user-defined retainedResults parameter,
// the oldest items are deleted.
// If retainedResults = 0, no metric is ever deleted.
func (w *Website) SaveResult(p *PollResult, retainedResults int) {
	itemsToDelete := 0
	if (retainedResults != 0) && (len(w.PollResults) >= retainedResults) {
		itemsToDelete = len(w.PollResults) + 1 - retainedResults
	} // If retainedResults is set to 0, store an unlimited number of results

	w.PollResults = append(w.PollResults[itemsToDelete:], *p)
}

// schedulePolls schedules regular polls for the website.
// The polling interval and the metrics retention policy are defined by the config file.
func (w *Website) schedulePolls(p PollConfig) {
	for range time.Tick(time.Duration(p.Interval) * time.Second) {
		w.Poll(p.RetainedResults)
	}
}

// Aggregate returns a payload.Metric containing the statistics for the website,
// aggregated over the specified timespan in seconds.
func (w *Website) Aggregate(timespan int) payload.Metric {
	// Copy poll results to ensure that they are not modified by
	// concurrent functions while results are being aggregated
	// TODO: avoid this somehow?
	p := w.PollResults
	startIdx := p.StartIndexFor(timespan)
	return payload.Metric{
		Availability:     p.Availability(startIdx),
		Average:          p.Average(startIdx),
		Max:              p.Max(startIdx),
		StatusCodeCounts: p.CountCodes(startIdx),
		ErrorCounts:      p.CountErrors(startIdx),
	}
}
