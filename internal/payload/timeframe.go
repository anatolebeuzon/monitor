package payload

import "time"

// Timeframe represents the time window that should be used
// to aggregate metrics.
type Timeframe struct {
	StartDate time.Time
	EndDate   time.Time
	Seconds   int
}

// NewTimeframe returns a new Timeframe, with the current date as the EndDate.
// The timespan input should be the duration, in seconds, between StartDate and EndDate.
func NewTimeframe(timespan int) Timeframe {
	ref := time.Now()
	return Timeframe{ref.Add(-time.Duration(timespan) * time.Second), ref, timespan}
}
