package client

import (
	"fmt"
	"go-project-3/types"

	"github.com/urfave/cli"
)

var Show = cli.Command{
	Name:  "show",
	Usage: "Show the dashboard",
	Action: func(c *cli.Context) error {
		// Load config
		var agg types.AggregateMetrics
		receivedData := make(chan bool)
		go schedulePolls(&agg, receivedData)
		fmt.Println("Fetching data...")
		Dashboard(&agg, receivedData)
		fmt.Println("Data fetched")
		return nil
	},
}

var Stop = cli.Command{
	Name:  "stop",
	Usage: "Stop daemon",
	Action: func(c *cli.Context) error {
		StopDaemon()
		return nil
	},
}
