package client

import (
	"monitor/payload"
	"sort"
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
	Bars    ui.BarChart
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

	bars := ui.NewBarChart()
	bars.BorderLabel = "Response code counts"
	bars.Height = 20
	bars.TextColor = ui.ColorGreen
	bars.BarColor = ui.ColorRed
	bars.NumColor = ui.ColorYellow

	return DashboardPage{
		Title:   *title,
		Counter: *counter,
		Metrics: *metrics,
		Alerts:  *alerts,
		Bars:    *bars,
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
		ui.NewRow(ui.NewCol(6, 0, &d.page.Alerts), ui.NewCol(6, 0, &d.page.Bars)),
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
	// test, _ := ExtractFromMap(s.Metrics[url][0].StatusCodeCounts)
	// p.Alerts.Text = fmt.Sprintf("%v", s.Metrics[url][s.Timespans.Order[0]].StatusCodeCounts)
	p.Bars.Data, p.Bars.DataLabels = GenerateBarChart(s.Metrics[url][s.Timespans.Order[0]])
	// p.Bars.Data, p.Bars.DataLabels = []int{16, 18, 13}, []string{"hi", "he", "ha"}
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

func GenerateBarChart(m payload.Metric) (values []int, labels []string) {
	var keys sort.IntSlice
	for key := range m.StatusCodeCounts {
		keys = append(keys, key)
	}
	keys.Sort()

	for _, key := range keys {
		values = append(values, m.StatusCodeCounts[key])
		labels = append(labels, strconv.Itoa(key))
	}

	labels = append(labels, "err")
	values = append(values, Count(m.ErrorCounts))
	return
}

func Count(errors map[string]int) (c int) {
	for _, i := range errors {
		c += i
	}
	return
}
