package agent

import (
	"fmt"
	"monitor/payload"
)

type Websites []Website

func NewWebsites(URLs []string) (w Websites) {
	for _, url := range URLs {
		w = append(w, Website{URL: url})
	}
	return
}

func (w Websites) SchedulePolls(p PollConfig) {
	for i := range w {
		go w[i].schedulePolls(p)
	}
	fmt.Println("All checks launched.")
}

// Alerts returns a payload.Alerts containing the alerts for all the websites.
//
// For each website, it compares its availability (on average, over the specified timespan in seconds)
// against the threshold, and creates an alert if the threshold is crossed.
func (w *Websites) Alerts(timespan int, threshold float64) payload.Alerts {
	alerts := make(payload.Alerts)
	for i, website := range *w {
		avail := website.Availability(timespan)
		if (avail < threshold) && !website.DownAlertSent {
			// if the website is considered down but no alert for this event was sent yet
			// create a "website is down" alert
			alerts[website.URL] = payload.NewAlert(avail, true)
			(*w)[i].DownAlertSent = true
		} else if (avail >= threshold) && website.DownAlertSent {
			// if the website is considered up but website was last reported down
			// create a "website has recovered" alert
			alerts[website.URL] = payload.NewAlert(avail, false)
			(*w)[i].DownAlertSent = false
		}
	}
	return alerts
}

// Stats returns a payload.Stats containing the aggregate statistics for all the websites.
func (w *Websites) Stats(timespan int) payload.Stats {
	p := payload.NewStats(timespan)
	for _, website := range *w {
		p.Metrics[website.URL] = website.Aggregate(timespan)
	}
	return p
}
