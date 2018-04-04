package client

import (
	"strconv"

	ui "github.com/gizak/termui"
)

// UIPage contains the UI objects used to display the dashboard.
//
// Contrary to UIDashboard, UIPage only contains elements that are visible:
// it does not store the data of websites that are not currently being shown.
type UIPage struct {
	Title   ui.Par // Shows the URL
	Counter ui.Par // Shows the index of the currently displayed website (e.g. 3/8)
	Left    UISide // Stats presented on the left-hand side of the dashboard
	Right   UISide // Stats presented on the right-hand side of the dashboard
	Alerts  ui.Par // Shows the latest alerts
}

// NewUIPage initializes the widgets of the dashboard with the
// appropriate UI parameters and returns a new DashboardPage.
func NewUIPage(c *Config) UIPage {
	Title := ui.NewPar("")
	Title.Height = 3

	Counter := ui.NewPar("")
	Counter.Height = 3
	Counter.Border = false

	Alerts := ui.NewPar("")
	Alerts.Height = 15
	Alerts.BorderLabel = "Alerts (aggregated over " + strconv.Itoa(c.Alerts.Timespan) + "s, "
	Alerts.BorderLabel += "refreshed every " + strconv.Itoa(c.Alerts.Frequency) + "s)"

	return UIPage{
		*Title,
		*Counter,
		NewUISide(c.Statistics.Left, ui.ColorBlue),
		NewUISide(c.Statistics.Right, ui.ColorYellow),
		*Alerts,
	}
}

// Refresh rerenders the DashboardPage using the latest data available.
func (p *UIPage) Refresh(s *Store) {
	s.RLock()
	defer s.RUnlock()

	url := s.URLs[s.currentIdx]

	// Update top-level widgets
	p.Title.Text = url
	p.Counter.Text = "Page " + strconv.Itoa(s.currentIdx+1) + "/" + strconv.Itoa(len(s.URLs))
	p.Alerts.Text = PrintAlert(&s.Alerts, url)

	// Update stats on both sides
	p.Left.Refresh(s.Metrics[url][p.Left.Timespan])
	p.Right.Refresh(s.Metrics[url][p.Right.Timespan])
}

func PrintAlert(a *Alerts, url string) (str string) {
	for _, alert := range (*a)[url] {
		str += "Website " + url + " is "
		if alert.BelowThreshold {
			str += "down. "
		} else {
			str += "up. "
		}
		str += "availability=" + strconv.FormatFloat(alert.Availability, 'f', 3, 64)
		str += ", time=" + alert.Timeframe.EndDate.String() + "\n"
	}
	return
}
