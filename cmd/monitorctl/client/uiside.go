/*
This file handles the display and update logic of UI elements presenting
statistics on the left and right frames ("sides") of the dashboard.
*/

package client

import (
	"sort"
	"strconv"
	"time"

	"github.com/oxlay/monitor/internal/payload"

	ui "github.com/gizak/termui"
)

// UISide contains the UI objects used to display the stats
// on either side of the dashboard.
type UISide struct {
	Timespan     int          // The timespan by which metrics are aggregated on this side
	Title        ui.Par       // Title of the aggregate
	Availability ui.Gauge     // Availability gauge
	Breakdown    ui.Table     // HTTP lifecycle steps durations
	CodeCounts   ui.BarChart  // Bar chart of the HTTP response codes counts
	RespGraph    ui.LineChart // Graph of response time evolution
	Errors       ui.Par       // Latest client (non-HTTP) errors
}

// NewUISide initializes the widgets of the dashboard side with the
// appropriate UI parameters and returns a new UISide.
func NewUISide(t TimeConf, color ui.Attribute) UISide {
	Title := ui.NewPar("")
	Title.Text = "Aggregated over " + strconv.Itoa(t.Timespan) + "s"
	Title.Text += " (refreshed every " + strconv.Itoa(t.Frequency) + "s)"
	Title.Height = 1
	Title.Border = false

	Availability := ui.NewGauge()
	Availability.BorderLabel = "Availability"
	Availability.Height = 3
	Availability.BorderFg = color

	Breakdown := ui.NewTable()
	Breakdown.BorderLabel = "Request breakdown"
	Breakdown.Rows = [][]string{
		[]string{"", "DNS", "TCP", "TLS", "Srv Process", "TTFB", "Transfer", "Response"},
		[]string{}, // average values ; will be populated during render
		[]string{}, // max values 	  ; same
	}
	Breakdown.FgColor = ui.ColorWhite
	Breakdown.BgColor = ui.ColorDefault
	Breakdown.BorderFg = color
	Breakdown.Height = 5
	Breakdown.TextAlign = ui.AlignCenter
	Breakdown.Separator = false

	CodeCounts := ui.NewBarChart()
	CodeCounts.BorderLabel = "Response code counts"
	CodeCounts.Height = 10
	CodeCounts.BorderFg = color

	RespGraph := ui.NewLineChart()
	RespGraph.BorderLabel = "Average response time evolution"
	RespGraph.Height = 10
	RespGraph.BorderFg = color
	RespGraph.Mode = "dot"
	RespGraph.DotStyle = '+'

	Errors := ui.NewPar("")
	Errors.BorderLabel = "Latest errors"
	Errors.Height = 7
	Errors.BorderFg = color

	return UISide{
		t.Timespan,
		*Title,
		*Availability,
		*Breakdown,
		*CodeCounts,
		*RespGraph,
		*Errors,
	}
}

// Refresh updates the UISide using the latest available data.
func (s *UISide) Refresh(m Metric) {
	// Update availability gauge
	s.Availability.Percent = int(m.Latest.Availability * 100)

	// Update color of the availability gauge
	avail := s.Availability.Percent
	if avail > 90 {
		s.Availability.BarColor = ui.ColorGreen
	} else if avail > 70 {
		s.Availability.BarColor = ui.ColorYellow
	} else {
		s.Availability.BarColor = ui.ColorRed
	}

	// Update request timing breakdown
	s.Breakdown.Rows[1] = FormatForTable("Avg", m.Latest.Average)
	s.Breakdown.Rows[2] = FormatForTable("Max", m.Latest.Max)

	// Update code counts
	s.CodeCounts.DataLabels, s.CodeCounts.Data = ExtractResponseCounts(m.Latest)

	// Update response time graph
	s.RespGraph.Data = FormatForGraph(m.AvgRespHist)

	// Update errors list
	s.Errors.Text = "" // Reset ErrorCounts text
	for err, c := range m.Latest.ErrorCounts {
		s.Errors.Text += err + " (" + strconv.Itoa(c) + " times)\n"
	}
}

// ExtractResponseCounts reads a metric and returns the corresponding
// slices that can be used to display a ui.BarChart of response code counts.
//
// For example, given:
//	m.StatusCodeCounts = map[int]int{200: 5, 404: 2, 500: 1}
//	m.ErrorCounts = map[string]int{"dial tcp: i/o timeout": 2}
// ExtractResponseCounts will return:
// 	codeNames  = []string{"200", "404", "5O0", "err"}
//	codeCounts = []int{5, 2, 1, 2}
func ExtractResponseCounts(m payload.Metric) (codeNames []string, codeCounts []int) {
	// Gather all the HTTP response codes and sort them in ascending order
	var codes sort.IntSlice
	for code := range m.StatusCodeCounts {
		codes = append(codes, code)
	}
	codes.Sort()

	// Generate code labels and code counts
	for _, code := range codes {
		codeNames = append(codeNames, strconv.Itoa(code))
		codeCounts = append(codeCounts, m.StatusCodeCounts[code])
	}

	// Append client (non-HTTP) error count at the end
	codeNames = append(codeNames, "err")
	codeCounts = append(codeCounts, Count(m.ErrorCounts))

	return
}

// Count returns the total number of errors in the input map.
func Count(errors map[string]int) (c int) {
	for _, i := range errors {
		c += i
	}
	return
}

// FormatForTable formats a timing for use in a ui.Table.
// Durations are rounded to the nearest millisecond.
//
// A sample return result would be:
// 	[]string{prefix, "12ms", "56ms", "87ms", ...}
func FormatForTable(prefix string, t payload.Timing) (s []string) {
	s = append(s, prefix)
	durations := []time.Duration{t.DNS, t.TCP, t.TLS, t.Server, t.TTFB, t.Transfer, t.Response}
	for _, d := range durations {
		s = append(s, d.Round(time.Millisecond).String())
	}
	return
}

// FormatForGraph formats a slice of durations for use in a ui.LineChart (i.e. a graph).
//
// Durations are rounded to the nearest millisecond, and converted to float64 values.
// float64(1) represents one second.
func FormatForGraph(d []time.Duration) (f []float64) {
	for _, duration := range d {
		f = append(f, float64(duration/time.Millisecond)/1000)
	}
	return
}
