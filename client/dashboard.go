package client

import (
	"monitor/payload"
	"sort"
	"strconv"
	"time"

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

func ExtractResponseCounts(m payload.Metric) (values []int, labels []string) {
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

func ToString(prefix string, d []time.Duration) (s []string) {
	s = append(s, prefix)
	for _, duration := range d {
		s = append(s, duration.Round(time.Millisecond).String())
	}
	return
}

func ToInt(d []time.Duration) (i []int) {
	for _, duration := range d {
		i = append(i, int(duration/time.Millisecond))
	}
	return
}

func ToFloat64(d []time.Duration) (f []float64) {
	for _, duration := range d {
		f = append(f, float64(duration/time.Millisecond)/1000)
	}
	return
}
