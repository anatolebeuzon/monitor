package client

import (
	"monitor/payload"
	"sort"
	"strconv"
	"time"

	ui "github.com/gizak/termui"
)

// DashboardSide contains the UI objects used to display the stats
// on either side of the dashboard.
type DashboardSide struct {
	Timespan     int          // The timespan by which metrics are aggregated
	Title        ui.Par       // Title of the aggregate
	Availability ui.Gauge     // Availability gauge
	Breakdown    ui.Table     // HTTP lifecycle steps durations
	CodeCounts   ui.BarChart  // Bar chart of the HTTP response codes counts
	RespGraph    ui.LineChart // Graph of response time evolution
	Errors       ui.Par       // Latest client (non-HTTP) errors
}

// NewDashboardSide initializes the widgets of the dashboard side with the
// appropriate UI parameters and returns a new DashboardPage.
func NewDashboardSide(t TimeConf, color ui.Attribute) DashboardSide {
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

	return DashboardSide{
		t.Timespan,
		*Title,
		*Availability,
		*Breakdown,
		*CodeCounts,
		*RespGraph,
		*Errors,
	}
}

// Refresh updates the DashboardSide using the latest data available.
func (s *DashboardSide) Refresh(m Metric) {
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
	s.Breakdown.Rows[1] = ToString("Avg", m.Latest.Average.ToSlice())
	s.Breakdown.Rows[2] = ToString("Max", m.Latest.Max.ToSlice())

	// Update code counts
	s.CodeCounts.Data, s.CodeCounts.DataLabels = ExtractResponseCounts(m.Latest)

	// Update response time graph
	s.RespGraph.Data = ToFloat64(m.AvgRespHist)

	// Update errors list
	s.Errors.Text = "" // Reset ErrorCounts text
	for err, c := range m.Latest.ErrorCounts {
		s.Errors.Text += err + " (" + strconv.Itoa(c) + " times)\n"
	}
}

// ExtractResponseCounts reads a metric and returns the corresponding
// slices that can be used to display a bar chart of response code counts.
//
// For example, the
func ExtractResponseCounts(m payload.Metric) (codeCounts []int, codeNames []string) {
	// Gather all the HTTP response codes and sort them in ascending order
	var codes sort.IntSlice
	for code := range m.StatusCodeCounts {
		codes = append(codes, code)
	}
	codes.Sort()

	// Generate code count and labels
	for _, code := range codes {
		codeCounts = append(codeCounts, m.StatusCodeCounts[code])
		codeNames = append(codeNames, strconv.Itoa(code))
	}

	// Append client (non-HTTP) error count at the end
	codeNames = append(labels, "err")
	codeCounts = append(codeCounts, Count(m.ErrorCounts))

	return
}

func Count(errors map[string]int) (c int) {
	for _, i := range errors {
		c += i
	}
	return
}

func ToString(prefix string, d []time.Duration) (s []string) {
	s = append(s, prefix)
	for _, duration := range d {
		s = append(s, duration.Round(time.Millisecond).String())
	}
	return
}

func ToFloat64(d []time.Duration) (f []float64) {
	for _, duration := range d {
		f = append(f, float64(duration/time.Millisecond)/1000)
	}
	return
}
