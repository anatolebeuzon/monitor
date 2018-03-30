package agent

import (
	"time"
)

func (w *Website) schedulePolls(pollInterval int) {
	for range time.Tick(time.Duration(pollInterval) * time.Second) {
		w.Poll()
	}
}

func (w *Websites) schedulePolls(pollInterval int) {
	for i := range *w {
		go (*w)[i].schedulePolls(pollInterval)
	}
}
