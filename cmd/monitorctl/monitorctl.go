/*
monitorctl is a client for the monitord daemon.
It displays website statistics and alerts on a console dashboard.

Usage :
	monitorctl [-config path]
where path is the relative path to the config file of the client.
If the config flag is not provided, monitorctl will look for
a file named config.json in the current directory.

Note that monitorctl's config file is different from monitord's.

Once the dashboard is shown, you can navigate between websites using left and
right arrows, or press "Q" to quit the dashboard.

Configuration

A sample JSON config file is described below:
	{
		"Server": "127.0.0.1:1234",	// Address on which monitord listens
		"Statistics": {
			"Left": {				// Left side of the dashboard
			"Frequency": 2,			// Frequency at which the daemon should be polled for stats
			"Timespan": 20			// Timespan over which metrics should be aggregated
			},
			"Right": {				// Right side of the dashboard
			"Frequency": 10,
			"Timespan": 40
			}
		},
		"Alerts": {
			"Frequency": 4,			// Frequency at which the daemon should be polled for alerts
			"Timespan": 120			// Timespan over which average availability should be computed
		}
	}
*/
package main

import (
	"flag"
	"monitor/client"
)

func main() {
	// Load config
	path := flag.String("config", "config.json", "Config file in JSON format")
	flag.Parse()
	config := client.ReadConfig(*path)

	// Create a new store to store the information received from the daemon
	store := client.NewStore()

	// Create new scheduler to regularly poll the daemon
	f := client.NewFetcher(config, store)
	f.Init() // start polling

	// Create and display a new dashboard
	d := client.NewUIDashboard(store, &config, f.UpdateUI)
	d.Show()
}
