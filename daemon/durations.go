package agent

import (
	"time"
)

// avgDuration returns the average duration of a slice of time.Duration.
func avgDuration(durations []time.Duration) time.Duration {
	avg := time.Duration(0)
	for _, duration := range durations {
		avg += duration
	}
	if len(durations) != 0 {
		avg /= time.Duration(len(durations))
	}
	return avg
}

// minDuration returns the minimal duration of a slice of time.Duration.
func minDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		// TODO: improve this error handling ?
		return time.Duration(0)
	}
	min := durations[0]
	for _, duration := range durations {
		if duration < min {
			min = duration
		}
	}
	return min
}

// maxDuration returns the maximal duration of a slice of time.Duration.
func maxDuration(durations []time.Duration) time.Duration {
	max := time.Duration(0)
	for _, duration := range durations {
		if duration > max {
			max = duration
		}
	}
	return max
}
