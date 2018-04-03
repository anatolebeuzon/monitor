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
	// Hostname string
	URL         string
	PollResults PollResults

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

	var t0, t1, t2, t3, t4, t5, t6 time.Time

	trace := &httptrace.ClientTrace{
		DNSStart:             func(_ httptrace.DNSStartInfo) { t0 = time.Now() },
		DNSDone:              func(_ httptrace.DNSDoneInfo) { t1 = time.Now() },
		ConnectStart:         func(_, _ string) { t2 = time.Now() },
		ConnectDone:          func(_, _ string, err error) { t3 = time.Now() }, // TODO: handle err?
		GotConn:              func(_ httptrace.GotConnInfo) { t4 = time.Now() },
		GotFirstResponseByte: func() { t5 = time.Now() },
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
		t6 = time.Now() // body has been read, response is over
	}

	// If an error occured, some times may still be at 0
	// However, they must be handled to avoid getting absurd results
	if t2.IsZero() {
		t2 = t1
	}
	if t3.IsZero() {
		t3 = t2
	}
	if t4.IsZero() {
		t4 = t3
	}
	if t5.IsZero() {
		t5 = t4
	}
	if t6.IsZero() {
		t6 = t5
	}

	p.Date = t0
	p.Timing = payload.Timing{
		DNS:      t1.Sub(t0),
		TCP:      t3.Sub(t2),
		TLS:      t4.Sub(t3),
		Server:   t5.Sub(t4),
		Transfer: t6.Sub(t5),
		TTFB:     t5.Sub(t0),
		Response: t6.Sub(t0),
	}

	fmt.Println(p)
	w.SaveResult(&p, retainedResults)
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

// Availability returns the average availability of a website over the specified timespan in seconds.
func (w *Website) Availability(timespan int) float64 {
	// TODO: remove duplicated code with aggregateResults /!\

	// Copy poll results to ensure that they are not modified by
	// concurrent functions while results are being aggregated
	p := w.PollResults
	startIdx := p.StartIndexFor(timespan)
	return p.Availability(startIdx)
}

// Aggregate returns a payload.Metric containing the statistics for the website,
// aggregated over the specified timespan in seconds.
func (w *Website) Aggregate(timespan int) payload.Metric {
	// Copy poll results to ensure that they are not modified by
	// concurrent functions while results are being aggregated
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
