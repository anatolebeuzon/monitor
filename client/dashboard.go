package client

import (
	"go-project-3/types"

	ui "github.com/gizak/termui"
)

type Dashboard struct {
	agg        *types.AggregateMapByURL
	currentIdx int
	page       DashboardPage
	updateUI   chan bool
}

type DashboardPage struct {
	Title   ui.Par
	Metrics ui.Par
	Alerts  ui.Par
}

func NewDashboardPage() DashboardPage {
	title := ui.NewPar("")
	title.Height = 3

	metrics := ui.NewPar("")
	metrics.Height = 20
	metrics.BorderLabel = "Metrics"

	alerts := ui.NewPar("")
	alerts.Height = 20
	alerts.BorderLabel = "Alerts"

	return DashboardPage{
		Title:   *title,
		Metrics: *metrics,
		Alerts:  *alerts,
	}
}

func NewDashboard(agg *types.AggregateMapByURL, updateUI chan bool) (d Dashboard) {
	return Dashboard{
		agg:        agg,
		currentIdx: 0,
		page:       NewDashboardPage(),
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
		ui.NewRow(ui.NewCol(2, 5, &d.page.Title)),
		ui.NewRow(ui.NewCol(12, 0, &d.page.Metrics)),
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
	d.page.Refresh(d.currentIdx, *d.agg)
	// Rerender UI
	ui.Render(ui.Body)
}

func (p *DashboardPage) Refresh(currentIdx int, agg types.AggregateMapByURL) {
	url := agg.URLs[currentIdx]
	p.Title.Text = url
	p.Metrics.Text = agg.String(url)
	p.Alerts.Text = "" // TODO
}

func (d *Dashboard) RegisterEventHandlers() {
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<right>", func(ui.Event) {
		if d.currentIdx < len((*d.agg).URLs)-1 {
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
