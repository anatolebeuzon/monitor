package agent

import (
	"fmt"

	"github.com/urfave/cli"
)

const path = "agent/config.json"

var Start = cli.Command{
	Name:  "start",
	Usage: "Start the agent",
	Action: func(c *cli.Context) error {
		// Load config
		config, err := readConfig(path)
		if err != nil {
			return err
		}
		websites := NewWebsites(config.URLs)
		fmt.Println(websites)

		websites.schedulePolls(config.PollInterval)
		websites.ServeRPC(config.ListeningPort)
		return nil
	},
}
