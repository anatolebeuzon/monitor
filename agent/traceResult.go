package agent

import (
	"time"
)

type TraceResults []TraceResult

type TraceResult struct {
	Date        time.Time
	DNStime     time.Duration
	TLStime     time.Duration
	ConnectTime time.Duration
	TTFB        time.Duration
	Error       error
	StatusCode  int
}

func (t TraceResults) startIndexFor(timespan int) int {
	threshold := time.Now().Add(-time.Duration(timespan) * time.Second)
	for i := len(t) - 1; i >= 0; i-- {
		if t[i].Date.Before(threshold) {
			return i + 1 // TODO: handle case where i + 1 is out of range
		}
	}
	return 0
}

func (t TraceResults) TTFBs(startIdx int) (durations []time.Duration) {
	for i := startIdx; i < len(t); i++ {
		durations = append(durations, t[i].TTFB)
	}
	return
}

func (t TraceResults) CountCodes(startIdx int) map[int]int {
	codesCount := make(map[int]int)
	for i := startIdx; i < len(t); i++ {
		codesCount[t[i].StatusCode]++
	}
	return codesCount
}

func (t TraceResults) Availability(startIdx int) float64 {
	c := 0
	for i := startIdx; i < len(t); i++ {
		if t[i].IsValid() {
			c++
		}
	}
	return float64(c) / float64(len(t)-startIdx)
}

func (t *TraceResult) IsValid() bool {
	return (t.Error == nil) && (t.StatusCode < 400)
}
