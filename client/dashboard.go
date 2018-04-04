package client

import (
	"log"

	ui "github.com/gizak/termui"
)

// Dashboard represents a dashboard, including the UI elements (widgets)
// and a pointer to the underlying data.
type Dashboard struct {
	store    *Store        // Pointer to the data that the dashboard can present
	page     DashboardPage // page contains the widgets that are presented on the dashboard
	updateUI chan bool     // updateUI signals when the dashboard should be rerendered (e.g. when new data arrives)
}

// NewDashboard returns a new dashboard.
func NewDashboard(s *Store, c *Config, updateUI chan bool) (d Dashboard) {
	return Dashboard{s, NewDashboardPage(c), updateUI}
}

// Show displays the dashboard on the console.
// It blocks until the user exits the dashboard.
func (d *Dashboard) Show() {
	// Initialize termui library
	if err := ui.Init(); err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// Handle user interactions
	d.RegisterEventHandlers()
	close := make(chan bool)
	go func() {
		ui.Loop()     // Handle keyboard events
		close <- true // If reached, it means ui.StopLoop() was called. Quit the dashboard
	}()

	d.BuildLayout()

	for {
		select {
		case <-d.updateUI:
			// Refresh the widgets with the latest data
			d.page.Refresh(d.store)

			// Rerender UI
			ui.Render(ui.Body)
		case <-close:
			// Quit the dashboard
			return
		}
	}
}

// BuildLayout builds the layout of widgets on the dashboard.
//
// It should only be called once.
// It does not need to be called again when dashboard data is updated.
func (d *Dashboard) BuildLayout() {
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(2, 5, &d.page.Title),    // Website URL
			ui.NewCol(1, 4, &d.page.Counter)), // Page counter
		ui.NewRow( // Aggregate titles
			ui.NewCol(4, 1, &d.page.Left.Title),
			ui.NewCol(4, 2, &d.page.Right.Title)),
		ui.NewRow( // Availability gauge and HTTP lifecycle detailed times
			ui.NewCol(6, 0, &d.page.Left.Availability, &d.page.Left.Breakdown),
			ui.NewCol(6, 0, &d.page.Right.Availability, &d.page.Right.Breakdown)),
		ui.NewRow(
			ui.NewCol(3, 0, &d.page.Left.CodeCounts), // Response code counts
			ui.NewCol(3, 0, &d.page.Left.RespGraph),  // Response time evolution graph
			ui.NewCol(3, 0, &d.page.Right.CodeCounts),
			ui.NewCol(3, 0, &d.page.Right.RespGraph),
		),
		ui.NewRow( // Latest client (non-HTTP) errors
			ui.NewCol(6, 0, &d.page.Left.Errors),
			ui.NewCol(6, 0, &d.page.Right.Errors),
		),
		ui.NewRow(ui.NewCol(12, 0, &d.page.Alerts)), // Latest alerts
	)
	ui.Body.Align() // Calculate layout based on window's width
}

// RegisterEventHandlers registers the keyboard
// events to which the dashboard will respond.
func (d *Dashboard) RegisterEventHandlers() {
	// Exit the dashboard when "Q" key is pressed
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	s := d.store
	s.Lock()
	defer s.Unlock()

	// Move to the next page when right arrow is pressed
	ui.Handle("/sys/kbd/<right>", func(ui.Event) {
		if s.currentIdx < len(s.URLs)-1 { // if there is a next page
			s.currentIdx++
			d.updateUI <- true
		}
	})

	// Move to the previous page when left arrow is pressed
	ui.Handle("/sys/kbd/<left>", func(ui.Event) {
		if s.currentIdx >= 1 { // if there is a previous page
			s.currentIdx--
			d.updateUI <- true
		}
	})
}
