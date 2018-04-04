package client

import (
	"time"
)

type Scheduler struct {
	Config Config
	Store  *Store

	// UpdateUI informs the dashboard that new data is available in the store,
	// and that it should re-render to display the latest information
	UpdateUI chan bool
}

// NewScheduler creates a new Scheduler with the provided Config.
func NewScheduler(c Config, s *Store) *Scheduler {
	return &Scheduler{
		Config:   c,
		Store:    s,
		UpdateUI: make(chan bool),
	}
}

// Init initiates regular polling of the daemon.
func (s *Scheduler) Init() {
	// Launch stat check routines
	c := &s.Config
	for _, t := range []TimeConf{c.Statistics.Left, c.Statistics.Right} {
		go func(t TimeConf) {
			s.GetStats(t.Timespan)
			for range time.Tick(time.Duration(t.Frequency) * time.Second) {
				s.GetStats(t.Timespan)
			}
		}(t)
	}

	// Launch alert check routine
	go func() {
		for range time.Tick(time.Duration(c.Alerts.Frequency) * time.Second) {
			s.GetAlerts(c.Alerts.Timespan)
		}
	}()
}
