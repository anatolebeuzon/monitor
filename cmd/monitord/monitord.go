/*
monitord is a daemon that polls websites, gathers related metrics,
and publishes those through an RPC API.

Usage :
	monitord [-config path]
where path is the relative path to the JSON config file of the daemon.
If the config flag is not provided, monitord will look for
a file named config.json in the current directory.

Note that monitord's config file is different from monitorctl's.

Configuration

A sample JSON config file is described below:

	{
		"ListeningPort": 1234, 			// the port on which the RPC server listens
		"Default": {
			"Interval": 2, 				// the interval, in seconds, between two requests to a given website
			"RetainedResults": 1000, 	// the number of poll results that are retained for a given website
			"Threshold": 0.8			// the availability threshold that triggers an alert when crossed
		},
		"Websites": [					// Websites to poll
			{
				"URL": "https://www.datadoghq.com",
				"Interval": 5,						// Defaults can be overridden on a per-website basis
				"RetainedResults": 5000,
				"Threshold": 0.95
			},
			{ "URL": "https://golang.org" }
  		]
	}
*/
package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/anatolebeuzon/monitor/cmd/monitord/daemon"
)

func main() {
	// Set up channel on which to send interrupt notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive the signal.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Load config
	path := flag.String("config", "", "Path to JSON config file")
	flag.Parse()
	config := daemon.ReadConfig(*path)

	// Start polling websites
	websites := daemon.NewWebsites(&config)
	websites.InitPolls()

	// Create RPC handler and start serving requests
	h := daemon.Handler(websites)
	daemon.ServeRPC(&h, config.ListeningPort, interrupt)

	return
}
