package client

import (
	"strconv"

	ui "github.com/gizak/termui"
)

type DashboardPage struct {
	Title   ui.Par
	Counter ui.Par
	Left    DashboardSide
	Right   DashboardSide
	Alerts  ui.Par
}

func NewDashboardPage(s *Store, c *Config) DashboardPage {
	Title := ui.NewPar("")
	Title.Height = 3

	Counter := ui.NewPar("")
	Counter.Height = 3
	Counter.Border = false

	Alerts := ui.NewPar("")
	Alerts.Height = 15
	Alerts.BorderLabel = "Alerts (aggregated over " + strconv.Itoa(c.Alerts.Timespan) + "s, "
	Alerts.BorderLabel += "refreshed every " + strconv.Itoa(c.Alerts.Frequency) + "s)"

	return DashboardPage{
		Title:   *Title,
		Counter: *Counter,
		Left:    NewDashboardSide(c.Statistics.Left, ui.ColorBlue),
		Right:   NewDashboardSide(c.Statistics.Right, ui.ColorYellow),
		Alerts:  *Alerts,
	}
}

func (p *DashboardPage) Refresh(currentIdx int, s Store) {
	url := s.URLs[currentIdx]
	p.Title.Text = url
	p.Counter.Text = "Page " + strconv.Itoa(currentIdx+1) + "/" + strconv.Itoa(len(s.URLs))
	p.Alerts.Text = s.Alerts.String(url)
	p.Left.Refresh(s.Metrics[url][p.Left.Timespan])
	p.Right.Refresh(s.Metrics[url][p.Right.Timespan])
}
