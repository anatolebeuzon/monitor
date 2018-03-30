package client

import (
	"go-project-3/types"

	ui "github.com/gizak/termui"
)

func Dashboard(agg *types.AggregateMetrics, receivedData chan bool) {

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	close := make(chan bool)
	go func() {
		ui.Loop()
		close <- true
	}()

	<-receivedData
	metric := (*agg)[0]

	title := ui.NewPar(metric.URL)
	title.Height = 3

	box0 := ui.NewPar(metric.String())
	box0.Height = 20
	box0.BorderLabel = "Metrics"

	alerts := ui.NewPar("")
	alerts.Height = 20
	alerts.BorderLabel = "Alerts"

	// build layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(2, 5, title),
		),
		ui.NewRow(
			ui.NewCol(12, 0, box0),
		),
		ui.NewRow(
			ui.NewCol(12, 0, alerts),
		),
	)

	// calculate layout
	ui.Body.Align()

	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	for {
		select {
		case <-receivedData:
			metric := (*agg)[0]
			box0.Text = metric.String()
			ui.Render(ui.Body)
		case <-close:
			return
		}

	}
}

func renderWebsite(website *Website) {

	// sinps := (func() []float64 {
	// 	n := 400
	// 	ps := make([]float64, n)
	// 	for i := range ps {
	// 		ps[i] = 1 + math.Sin(float64(i)/5)
	// 	}
	// 	return ps
	// })()
	// sinpsint := (func() []int {
	// 	ps := make([]int, len(sinps))
	// 	for i, v := range sinps {
	// 		ps[i] = int(100*v + 10)
	// 	}
	// 	return ps
	// })()

	// spark := ui.Sparkline{}
	// spark.Height = 8
	// spdata := sinpsint
	// spark.Data = spdata[:100]
	// spark.LineColor = ui.ColorCyan
	// spark.TitleColor = ui.ColorWhite

	// sp := ui.NewSparklines(spark)
	// sp.Height = 11
	// sp.BorderLabel = "Sparkline"

	// lc := ui.NewLineChart()
	// lc.BorderLabel = "braille-mode Line Chart"
	// lc.Data = sinps
	// lc.Height = 11
	// lc.AxesColor = ui.ColorWhite
	// lc.LineColor = ui.ColorYellow | ui.AttrBold

	// gs := make([]*ui.Gauge, 3)
	// for i := range gs {
	// 	gs[i] = ui.NewGauge()
	// 	//gs[i].LabelAlign = ui.AlignCenter
	// 	gs[i].Height = 2
	// 	gs[i].Border = false
	// 	gs[i].Percent = i * 10
	// 	gs[i].PaddingBottom = 1
	// 	gs[i].BarColor = ui.ColorRed
	// }

	// ls := ui.NewList()
	// ls.Border = false
	// ls.Items = []string{
	// 	"[1] Downloading File 1",
	// 	"", // == \newline
	// 	"[2] Downloading File 2",
	// 	"",
	// 	"[3] Uploading File 3",
	// }
	// ls.Height = 5

	// par := ui.NewPar("<> This row has 3 columns\n<- Widgets can be stacked up like left side\n<- Stacked widgets are treated as a single widget")
	// par.Height = 5
	// par.BorderLabel = "Demonstration"

	// ui.Handle("/timer/1s", func(e ui.Event) {
	// 	t := e.Data.(ui.EvtTimer)
	// 	i := t.Count
	// 	if i > 103 {
	// 		ui.StopLoop()
	// 		return
	// 	}

	// 	for _, g := range gs {
	// 		g.Percent = (g.Percent + 3) % 100
	// 	}

	// 	sp.Lines[0].Data = spdata[:100+i]
	// 	lc.Data = sinps[2*i:]
	// 	ui.Render(ui.Body)
	// })

	// ui.Handle("/sys/wnd/resize", func(e ui.Event) {
	// 	ui.Body.Width = ui.TermWidth()
	// 	ui.Body.Align()
	// 	ui.Clear()
	// 	ui.Render(ui.Body)
	// })
}
