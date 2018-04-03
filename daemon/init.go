/*
OK

This file contains the main data types used by the daemon,
and the init logic used on daemon startup:
- to create Website objects from URLs
- to launch websites' poll schedulers
*/

package daemon

import (
	"fmt"
	"monitor/payload"
	"time"
)

// Websites represents all the websites to be polled.
type Websites []Website

// A Website object contains the identity (URL) of the website,
// as well as all the corresponding poll results.
type Website struct {
	URL         string
	PollResults PollResults

	// DownAlertSent is true if at the last alert check from the front-end,
	// the aggregate availability was below the threshold. Keeping this information:
	// - avoids sending repetitive "website is down!" alerts
	// - enables the sending of one "website is up!" alert upon website recovery
	DownAlertSent bool
}

// PollResults represents all the trace results for a given website.
//
// The retention policy of those results is user-defined: in the config file,
// the RetainedResults parameter specifies how many poll results to keep.
type PollResults []PollResult

// A PollResult represents the results of one request to a website.
type PollResult struct {
	// Date is the date at which the first byte of the response was received.
	Date time.Time

	// Timing contains the duration of the different phases of the request.
	Timing payload.Timing

	// Error stores the error if the request resulted in a client error, or nil otherwise.
	Error error

	// StatusCode stores the HTTP response code of the request, or 0 if the request
	// resulted in a (non-HTTP) client error.
	StatusCode int
}

// NewWebsites creates a new Websites object from a slice of URLs.
//
// NB: different URLs of the same domain (purposefully) lead
// to the creation of multiple Website objects.
func NewWebsites(URLs []string) (w Websites) {
	for _, url := range URLs {
		w = append(w, Website{URL: url})
	}
	return
}

// InitPolls launches, for each website, a poll scheduler in a separate goroutine.
func (w Websites) InitPolls(p PollConfig) {
	for i := range w {
		go w[i].SchedulePolls(p)
	}
	fmt.Println("All checks launched.")
}
