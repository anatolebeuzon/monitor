package agent

import (
	"time"
)

const pollingInterval = 2 // in seconds

func (website *Website) schedulePolls(pollInterval int) {
	for range time.Tick(time.Duration(pollInterval) * time.Second) {
		website.Poll()
	}
}

func (websites *Websites) schedulePolls(pollInterval int) {
	for i := range *websites {
		go (*websites)[i].schedulePolls(pollInterval)
	}
}
