/*
monitord is a daemon that polls websites and gather related metrics.
Those metrics are accessible through an RPC API.

Usage :
	monitord [-config path]
where path is the relative path to the JSON config file of the daemon.
If the config flag is not provided, monitord will look for
a file named config.json in the current directory.

Note that monitord's config file is different from monitorctl's.

Sample config file

	{
		"ListeningPort": 1234, 			// the port on which the RPC server listens
		"Poll": {
			"Interval": 2,				// the interval, in seconds, between two requests to a given website
			"RetainedResults": 1000000  // the number of poll results that are retained for a given website
	},
		"AlertThreshold": 0.8,			// the availability threshold that triggers an alert when crossed
		"URLs": [						// the array of URLs of 
			"https://youtube.com",
			"https://www.youtube.com",
			"https://apple.com",
			"https://www.datadoghq.com"
		]
	}
*/
package main

import (
	"flag"
	"fmt"
	"monitor/daemon"
	"os"
	"os/signal"
)

const (
	name        = "monitord"
	description = "Daemon that polls websites and gather related metrics"
)

func main() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive the signal.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Load config
	configPath := flag.String("config", "config.json", "Config file in JSON format")
	flag.Parse()
	config := agent.ReadConfig(*configPath)

	websites := agent.NewWebsites(config.URLs)
	websites.SchedulePolls(config.Poll)

	// Create RPC handler
	h := &agent.Handler{
		Websites:       &websites,
		AlertThreshold: config.AlertThreshold,
	}
	agent.ServeRPC(h, config.ListeningPort, interrupt)

	// Handle interrupt by system signal
	// TODO: improve closing logic?
	fmt.Println("Closing properly...")
	return
}
