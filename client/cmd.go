package client

import (
	"github.com/urfave/cli"
)

const path = "client/config.json"

var Show = cli.Command{
	Name:  "show",
	Usage: "Show the dashboard",
	Action: func(c *cli.Context) error {
		// Load config
		config, err := readConfig(path)
		if err != nil {
			return err
		}

		s := newScheduler(config)
		agg := s.init()
		d := NewDashboard(agg, s.updateUI)
		d.Show()
		return nil
	},
}

var Stop = cli.Command{
	Name:  "stop",
	Usage: "Stop daemon",
	Action: func(c *cli.Context) error {
		// Load config
		config, err := readConfig(path)
		if err != nil {
			return err
		}

		StopDaemon(config.Server)
		return nil
	},
}
