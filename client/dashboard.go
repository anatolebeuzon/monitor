package client

import (
	"strconv"

	ui "github.com/gizak/termui"
)

type Dashboard struct {
	store      *Store
	currentIdx int
	page       DashboardPage
	updateUI   chan bool
}

type DashboardPage struct {
	Title   ui.Par
	Counter ui.Par
	Metrics ui.Par
	Alerts  ui.Par
}

func NewDashboardPage() DashboardPage {
	title := ui.NewPar("")
	title.Height = 3

	counter := ui.NewPar("")
	counter.Height = 3
	counter.Border = false

	metrics := ui.NewPar("")
	metrics.Height = 20
	metrics.BorderLabel = "Metrics"

	alerts := ui.NewPar("")
	alerts.Height = 20
	alerts.BorderLabel = "Alerts"

	return DashboardPage{
		Title:   *title,
		Counter: *counter,
		Metrics: *metrics,
		Alerts:  *alerts,
	}
}

func NewDashboard(s *Store, updateUI chan bool) (d Dashboard) {
	return Dashboard{
		store:      s,
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
		ui.NewRow(ui.NewCol(2, 5, &d.page.Title), ui.NewCol(1, 4, &d.page.Counter)),
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
	d.page.Refresh(d.currentIdx, *d.store)
	// Rerender UI
	ui.Render(ui.Body)
}

func (p *DashboardPage) Refresh(currentIdx int, s Store) {
	url := s.URLs[currentIdx]
	p.Title.Text = url
	p.Counter.Text = strconv.Itoa(currentIdx+1) + "/" + strconv.Itoa(len(s.URLs))
	p.Metrics.Text = s.Metrics.String(url, s.Timespans.Order)
	p.Alerts.Text = s.Alerts.String(url)
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
