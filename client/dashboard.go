package client

import (
	ui "github.com/gizak/termui"
)

type Dashboard struct {
	store      *Store
	currentIdx int
	page       DashboardPage
	updateUI   chan bool
}

func NewDashboard(s *Store, c *Config, updateUI chan bool) (d Dashboard) {
	return Dashboard{
		store:      s,
		currentIdx: 0,
		page:       NewDashboardPage(s, c),
		updateUI:   updateUI,
	}
}

func (d *Dashboard) Show() error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	close := make(chan bool)
	go func() {
		ui.Loop()
		close <- true
	}()

	d.RegisterEventHandlers()

	// build layout
	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(2, 5, &d.page.Title), ui.NewCol(1, 4, &d.page.Counter)),
		ui.NewRow(ui.NewCol(4, 1, &d.page.Left.Title), ui.NewCol(4, 2, &d.page.Right.Title)),
		ui.NewRow(
			ui.NewCol(6, 0, &d.page.Left.Availability, &d.page.Left.Breakdown),
			ui.NewCol(6, 0, &d.page.Right.Availability, &d.page.Right.Breakdown)),
		ui.NewRow(
			ui.NewCol(3, 0, &d.page.Left.CodeCounts),
			ui.NewCol(3, 0, &d.page.Left.RespHist),
			ui.NewCol(3, 0, &d.page.Right.CodeCounts),
			ui.NewCol(3, 0, &d.page.Right.RespHist),
		),
		ui.NewRow(
			ui.NewCol(6, 0, &d.page.Left.Errors),
			ui.NewCol(6, 0, &d.page.Right.Errors),
		),
		ui.NewRow(ui.NewCol(12, 0, &d.page.Alerts)),
	)
	ui.Body.Align()

	for {
		select {
		case <-d.updateUI:
			d.Render()
		case <-close:
			return nil
		}
	}
}

func (d *Dashboard) Render() {
	// Refresh the DashboardPage object currently used
	d.page.Refresh(d.currentIdx, *d.store)
	// Rerender UI
	ui.Render(ui.Body)
}

func (d *Dashboard) RegisterEventHandlers() {
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<right>", func(ui.Event) {
		if d.currentIdx < len(d.store.URLs)-1 {
			d.currentIdx++
			d.updateUI <- true
		}
	})

	ui.Handle("/sys/kbd/<left>", func(ui.Event) {
		if d.currentIdx >= 1 {
			d.currentIdx--
			d.updateUI <- true
		}
	})
}
