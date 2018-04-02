package client

import (
	"strconv"

	ui "github.com/gizak/termui"
)

type DashboardSide struct {
	Timespan     int
	Title        ui.Par
	Availability ui.Gauge
	Breakdown    ui.Table
	CodeCounts   ui.BarChart
	RespHist     ui.LineChart
	Errors       ui.Par
}

func NewDashboardSide(s Statistic, color ui.Attribute) DashboardSide {
	text := "Aggregated over " + strconv.Itoa(s.Timespan) + "s"
	text += " (refreshed every " + strconv.Itoa(s.Frequency) + "s)"
	Title := ui.NewPar(text)
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
		[]string{},
		[]string{},
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

	RespHist := ui.NewLineChart()
	RespHist.BorderLabel = "Average response time evolution"
	RespHist.Height = 10
	RespHist.BorderFg = color
	RespHist.Mode = "dot"
	RespHist.DotStyle = '+'

	Errors := ui.NewPar("")
	Errors.BorderLabel = "Latest errors"
	Errors.Height = 7
	Errors.BorderFg = color

	return DashboardSide{
		Timespan:     s.Timespan,
		Title:        *Title,
		Availability: *Availability,
		Breakdown:    *Breakdown,
		CodeCounts:   *CodeCounts,
		RespHist:     *RespHist,
		Errors:       *Errors,
	}
}

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

	// Update response history
	s.RespHist.Data = ToFloat64(m.AvgRespHist)

	// Update errors list
	s.Errors.Text = "" // Reset ErrorCounts text
	for err, c := range m.Latest.ErrorCounts {
		s.Errors.Text += err + " (" + strconv.Itoa(c) + " times)\n"
	}
}
