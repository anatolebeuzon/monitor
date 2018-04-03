package client

import (
	"strconv"

	ui "github.com/gizak/termui"
)

// DashboardPage contains the UI objects used to display the dashboard.
//
// Contrary to Dashboard, DashboardPage only contains elements that are visible:
// it does not store the data of websites that are not currently being shown.
type DashboardPage struct {
	Title   ui.Par        // Shows the URL
	Counter ui.Par        // Shows the index of the currently displayed website (e.g. 3/8)
	Left    DashboardSide // Stats presented on the left-hand side of the dashboard
	Right   DashboardSide // Stats presented on the right-hand side of the dashboard
	Alerts  ui.Par        // Shows the latest alerts
}

// NewDashboardPage initializes the widgets of the dashboard with the
// appropriate UI parameters and returns a new DashboardPage.
func NewDashboardPage(c *Config) DashboardPage {
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
		*Title,
		*Counter,
		NewDashboardSide(c.Statistics.Left, ui.ColorBlue),
		NewDashboardSide(c.Statistics.Right, ui.ColorYellow),
		*Alerts,
	}
}

// Refresh rerenders the DashboardPage using the latest data available.
func (p *DashboardPage) Refresh(currentIdx int, s Store) {
	url := s.URLs[currentIdx]

	// Update top-level widgets
	p.Title.Text = url
	p.Counter.Text = "Page " + strconv.Itoa(currentIdx+1) + "/" + strconv.Itoa(len(s.URLs))
	p.Alerts.Text = s.Alerts.String(url)

	// Update stats on both sides
	p.Left.Refresh(s.Metrics[url][p.Left.Timespan])
	p.Right.Refresh(s.Metrics[url][p.Right.Timespan])

	// Rerender UI
	ui.Render(ui.Body)
}
