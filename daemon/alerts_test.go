package daemon

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/oxlay/monitor/payload"
)

const testURL = "http://test/"

// End-to-end test of alert logic
// TODO: explain that I chose not to mock method calls
func TestAlerts(t *testing.T) {
	// Integration test - skip if necessary
	if testing.Short() {
		t.Skip()
	}

	end := time.Now()
	start := end.Add(-20 * time.Second)
	timeframe := payload.Timeframe{
		StartDate: start,
		EndDate:   end,
		Seconds:   20,
	}
	failure := PollResult{Date: end.Add(-1 * time.Second), StatusCode: 404}
	success := PollResult{Date: end.Add(-1 * time.Second), StatusCode: 200}

	// Create poll results at special dates regarding the timeframe created above
	beforeStartSuccess := PollResult{Date: start.Add(-1 * time.Second), StatusCode: 200}
	edgeStartSuccess := PollResult{Date: start, StatusCode: 200}
	edgeEndSuccess := PollResult{Date: end, StatusCode: 200}
	afterEndSuccess := PollResult{Date: end.Add(1 * time.Second), StatusCode: 200}

	testCases := []struct {
		handler  Handler
		expected payload.Alerts
	}{
		{
			// Website was up and is now down: alert expected
			buildHandler(false, failure, success),
			buildAlerts(timeframe, 0.5, true),
		},
		{
			// Website was down and is now up: alert expected
			buildHandler(true, success),
			buildAlerts(timeframe, 1, false),
		},
		{
			// Website was up and is still up: no alert expected
			buildHandler(false, success),
			payload.Alerts{},
		},
		{
			// Website was down and is still down: no alert expected
			buildHandler(true, failure),
			payload.Alerts{},
		},
		{
			// Website has no poll result available: alert expected
			buildHandler(false),
			buildAlerts(timeframe, 0, true),
		},
		{
			// Website has an old poll result: should be ignored in the availability calculation
			buildHandler(false, beforeStartSuccess, success, failure),
			buildAlerts(timeframe, 0.5, true),
		},
		{
			// Website has a poll result at the start edge of the timeframe:
			// should be included in the availability calculation
			buildHandler(false, edgeStartSuccess, failure),
			buildAlerts(timeframe, 0.5, true),
		},
		{
			// Website has a poll result at the end edge of the timeframe:
			// should be ignored in the availability calculation
			buildHandler(false, success, failure, edgeEndSuccess),
			buildAlerts(timeframe, 0.5, true),
		},
		{
			// Website has a poll result that is newer than the end of the timeframe:
			// should be ignored in the availability calculation
			buildHandler(false, success, failure, afterEndSuccess),
			buildAlerts(timeframe, 0.5, true),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint("Test case ", i), func(t *testing.T) {
			var computed payload.Alerts
			tc.handler.Alerts(timeframe, &computed)
			if !reflect.DeepEqual(computed, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, computed)
			}
		})
	}
}

func buildHandler(DownAlertSent bool, r ...PollResult) Handler {
	return Handler([]Website{Website{
		URL:           testURL,
		Threshold:     0.8,
		PollResults:   &PollResults{items: r},
		DownAlertSent: DownAlertSent,
	}})
}

func buildAlerts(tf payload.Timeframe, avail float64, belowThreshold bool) payload.Alerts {
	return payload.Alerts{
		testURL: payload.Alert{
			Timeframe:      tf,
			Availability:   avail,
			BelowThreshold: belowThreshold,
		},
	}
}
