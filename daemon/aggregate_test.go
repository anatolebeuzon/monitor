package daemon

// func TestExtract(t *testing.T) {
// 	end := time.Now()
// 	start := end.Add(-20 * time.Second)
// 	timeframe := payload.Timeframe{
// 		StartDate: start,
// 		EndDate:   end,
// 		Seconds:   20,
// 	}
// 	failure := PollResult{Date: end.Add(-1 * time.Second), StatusCode: 404}
// 	success := PollResult{Date: end.Add(-1 * time.Second), StatusCode: 200}

// 	// Create poll results at special dates regarding the timeframe created above
// 	beforeStartSuccess := PollResult{Date: start.Add(-1 * time.Second), StatusCode: 200}
// 	edgeStartSuccess := PollResult{Date: start, StatusCode: 200}
// 	edgeEndSuccess := PollResult{Date: end, StatusCode: 200}
// 	afterEndSuccess := PollResult{Date: end.Add(1 * time.Second), StatusCode: 200}

// 	testCases := struct{
// 		pollResults PollResults
// 		timespan int
// 		expected int
// 	}{
// 		[]PollResult{
// 			PollResult{Date: ref.Add(-21 * time.Second)},
// 			PollResult{Date: ref.Add(-15 * time.Second)},
// 		PollResult{Date: ref.Add(-21 * time.Second)},}
// 	}
// }
