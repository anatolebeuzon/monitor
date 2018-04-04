package payload

import "time"

// Timeframe represents the time window that should be used
// to aggregate metrics.
type Timeframe struct {
	StartDate time.Time
	EndDate   time.Time
	Seconds   int
}

func NewTimeframe(timespan int) Timeframe {
	ref := time.Now()
	return Timeframe{ref.Add(-time.Duration(timespan) * time.Second), ref, timespan}
}
