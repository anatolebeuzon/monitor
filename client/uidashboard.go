/*
This file contains the main logic of the dashboard, including:
- the init logic of the dashboard
- the layout computation
- the handling of dashboard events (such as keyboard events or window resize)
- the exit logic
*/

package client

import (
	"log"

	ui "github.com/gizak/termui"
)

// UIDashboard represents a dashboard, including the UI elements (widgets)
// and a pointer to the underlying data.
type UIDashboard struct {
	Store    *Store    // Pointer to the data that the dashboard can present
	Page     UIPage    // page contains the widgets that are presented on the dashboard
	UpdateUI chan bool // updateUI signals when the dashboard should be rerendered (e.g. when new data arrives)
}

// NewUIDashboard returns a new dashboard.
func NewUIDashboard(s *Store, c *Config, updateUI chan bool) UIDashboard {
	return UIDashboard{s, NewUIPage(c), updateUI}
}

// Show displays the dashboard on the console.
// It blocks until the user exits the dashboard.
func (d *UIDashboard) Show() {
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
		case <-d.UpdateUI:
			// Refresh the widgets with the latest data
			d.Page.Refresh(d.Store)

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
func (d *UIDashboard) BuildLayout() {
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(2, 5, &d.Page.Title),    // Website URL
			ui.NewCol(1, 4, &d.Page.Counter)), // Page counter
		ui.NewRow( // Aggregate titles
			ui.NewCol(4, 1, &d.Page.Left.Title),
			ui.NewCol(4, 2, &d.Page.Right.Title)),
		ui.NewRow( // Availability gauge and HTTP lifecycle detailed times
			ui.NewCol(6, 0, &d.Page.Left.Availability, &d.Page.Left.Breakdown),
			ui.NewCol(6, 0, &d.Page.Right.Availability, &d.Page.Right.Breakdown)),
		ui.NewRow(
			ui.NewCol(3, 0, &d.Page.Left.CodeCounts), // Response code counts
			ui.NewCol(3, 0, &d.Page.Left.RespGraph),  // Response time evolution graph
			ui.NewCol(3, 0, &d.Page.Right.CodeCounts),
			ui.NewCol(3, 0, &d.Page.Right.RespGraph),
		),
		ui.NewRow( // Latest client (non-HTTP) errors
			ui.NewCol(6, 0, &d.Page.Left.Errors),
			ui.NewCol(6, 0, &d.Page.Right.Errors),
		),
		ui.NewRow(ui.NewCol(12, 0, &d.Page.Alerts)), // Latest alerts
		ui.NewRow(ui.NewCol(12, 0, &d.Page.Footer)), // Footer with navigation information
	)
	ui.Body.Align() // Calculate layout based on window's width
}

// RegisterEventHandlers registers the keyboard
// events to which the dashboard will respond.
func (d *UIDashboard) RegisterEventHandlers() {
	// Exit the dashboard when "Q" key is pressed
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	// Handle window resize
	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	// Move to the next page when right arrow is pressed
	ui.Handle("/sys/kbd/<right>", func(ui.Event) {
		s := d.Store
		s.Lock()
		defer s.Unlock()

		if s.CurrentIdx < len(s.URLs)-1 { // if there is a next page
			s.CurrentIdx++
			d.UpdateUI <- true
		}
	})

	// Move to the previous page when left arrow is pressed
	ui.Handle("/sys/kbd/<left>", func(ui.Event) {
		s := d.Store
		s.Lock()
		defer s.Unlock()

		if s.CurrentIdx >= 1 { // if there is a previous page
			s.CurrentIdx--
			d.UpdateUI <- true
		}
	})
}
