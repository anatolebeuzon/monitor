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

type DashboardPage struct {
	Title   ui.Par
	Counter ui.Par
	Left    DashboardSide
	Right   DashboardSide
	Alerts  ui.Par
}

type DashboardSide struct {
	Timespan     int
	Title        ui.Par
	Availability ui.Gauge
	TTFB         ui.Par
	CodeCounts   ui.BarChart
	Errors       ui.Par
	Breakdown    ui.Table
}

func NewDashboardPage(s *Store, c *Config) DashboardPage {
	title := ui.NewPar("")
	title.Height = 3

	counter := ui.NewPar("")
	counter.Height = 3
	counter.Border = false

	alerts := ui.NewPar("")
	alerts.Height = 15
	alerts.BorderLabel = "Alerts (refreshed every " + strconv.Itoa(c.Alerts.Frequency) + "s)"

	return DashboardPage{
		Title:   *title,
		Counter: *counter,
		Left:    NewDashboardSide(c.Statistics.Left, ui.ColorBlue),
		Right:   NewDashboardSide(c.Statistics.Right, ui.ColorYellow),
		Alerts:  *alerts,
	}
}

func NewDashboardSide(s Statistic, color ui.Attribute) DashboardSide {
	text := "Aggregate over " + strconv.Itoa(s.Timespan) + "s"
	text += " (refreshed every " + strconv.Itoa(s.Frequency) + "s)"
	Title := ui.NewPar(text)
	Title.Height = 1
	Title.Border = false

	Availability := ui.NewGauge()
	Availability.BorderLabel = "Availability"
	Availability.Height = 3
	Availability.BorderFg = color

	TTFB := ui.NewPar("")
	TTFB.BorderLabel = "TTFBs"
	TTFB.Height = 5
	TTFB.BorderFg = color

	CodeCounts := ui.NewBarChart()
	CodeCounts.BorderLabel = "Response code counts"
	CodeCounts.Height = 25
	CodeCounts.BorderFg = color

	Errors := ui.NewPar("")
	Errors.BorderLabel = "Latest errors"
	Errors.Height = 10
	Errors.BorderFg = color

	Breakdown := ui.NewTable()
	Breakdown.BorderLabel = "Request breakdown"
	Breakdown.Rows = [][]string{
		[]string{"", "DNS lookup", "TCP connection", "TLS handshake", "Server processing", "Total (TTFB)"},
		[]string{},
		[]string{},
	}
	Breakdown.FgColor = ui.ColorWhite
	Breakdown.BgColor = ui.ColorDefault
	Breakdown.BorderFg = color
	Breakdown.Height = 10
	Breakdown.TextAlign = ui.AlignCenter
	Breakdown.Separator = false

	return DashboardSide{
		Timespan:     s.Timespan,
		Title:        *Title,
		Availability: *Availability,
		TTFB:         *TTFB,
		CodeCounts:   *CodeCounts,
		Errors:       *Errors,
		Breakdown:    *Breakdown,
	}
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
		ui.NewRow(ui.NewCol(6, 0, &d.page.Left.Availability), ui.NewCol(6, 0, &d.page.Right.Availability)),
		ui.NewRow(
			ui.NewCol(6, 0, &d.page.Left.TTFB, &d.page.Left.Breakdown, &d.page.Left.Errors),
			ui.NewCol(6, 0, &d.page.Right.TTFB, &d.page.Right.Breakdown, &d.page.Right.Errors),
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

func (p *DashboardPage) Refresh(currentIdx int, s Store) {
	url := s.URLs[currentIdx]
	p.Title.Text = url
	p.Counter.Text = "Page " + strconv.Itoa(currentIdx+1) + "/" + strconv.Itoa(len(s.URLs))
	p.Alerts.Text = s.Alerts.String(url)
	p.Left.Refresh(s.Metrics[url][p.Left.Timespan])
	p.Right.Refresh(s.Metrics[url][p.Right.Timespan])
}

func (s *DashboardSide) Refresh(m payload.Metric) {

	// Update availability gauge
	s.Availability.Percent = int(m.Availability * 100)

	// Update color of the availability gauge
	avail := s.Availability.Percent
	if avail > 90 {
		s.Availability.BarColor = ui.ColorGreen
	} else if avail > 80 {
		s.Availability.BarColor = ui.ColorYellow
	} else {
		s.Availability.BarColor = ui.ColorRed
	}

	// Update errors list
	s.Errors.Text = "" // Reset ErrorCounts text
	for err, c := range m.ErrorCounts {
		s.Errors.Text += err + " (" + strconv.Itoa(c) + " times)\n"
	}

	// Update TTFB stats
	s.TTFB.Text = "Average: " + m.Average.TTFB.String()

	// Update request timing breakdown
	s.Breakdown.Rows[1] = ToString("Avg", m.Average.ToSlice())
	s.Breakdown.Rows[2] = ToString("Max", m.Max.ToSlice())

	s.CodeCounts.Data, s.CodeCounts.DataLabels = ExtractResponseCounts(m)
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
