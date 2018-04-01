package agent

import (
	"testing"
	"time"
)

var test = struct {
	durations []time.Duration
	min       time.Duration
	max       time.Duration
	avg       time.Duration
}{
	[]time.Duration{time.Duration(1), time.Duration(2), time.Duration(3), time.Duration(4), time.Duration(5)},
	time.Duration(1),
	time.Duration(5),
	time.Duration(3),
}

func TestAvgDuration(t *testing.T) {
	actual := avgDuration(test.durations)
	if actual != test.avg {
		t.Errorf("avgDuration(%v): expected %v, got %v", test.durations, test.avg, actual)
	}
}

func TestMinDuration(t *testing.T) {
	actual := minDuration(test.durations)
	if actual != test.min {
		t.Errorf("minDuration(%v): expected %v, got %v", test.durations, test.min, actual)
	}
}

func TestMaxDuration(t *testing.T) {
	actual := maxDuration(test.durations)
	if actual != test.max {
		t.Errorf("maxDuration(%v): expected %v, got %v", test.durations, test.max, actual)
	}
}
