package agent

import (
	"fmt"
	"io/ioutil"
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

	var t0, t1, t2, t3, t4, t5, t6 time.Time

	trace := &httptrace.ClientTrace{
		DNSStart:             func(_ httptrace.DNSStartInfo) { t0 = time.Now() },
		DNSDone:              func(_ httptrace.DNSDoneInfo) { t1 = time.Now() },
		ConnectStart:         func(_, _ string) { t2 = time.Now() },
		ConnectDone:          func(_, _ string, err error) { t3 = time.Now() }, // TODO: handle err?
		GotConn:              func(_ httptrace.GotConnInfo) { t4 = time.Now() },
		GotFirstResponseByte: func() { t5 = time.Now() },
	}

	var tr TraceResult
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := NewTransport().RoundTrip(req)
	if err != nil {
		tr.Error = err
	} else {
		tr.StatusCode = resp.StatusCode
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	t6 = time.Now() // body has been read, response is over

	tr.Date = t5
	tr.Timing.DNS = t1.Sub(t0)
	tr.Timing.TCP = t3.Sub(t2)
	tr.Timing.TLS = t4.Sub(t3)
	tr.Timing.Server = t5.Sub(t4)
	tr.Timing.Transfer = t6.Sub(t5)
	tr.Timing.TTFB = t5.Sub(t0)
	tr.Timing.Response = t6.Sub(t0)

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
