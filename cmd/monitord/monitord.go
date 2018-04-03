/*
monitord is a daemon that polls websites and gather related metrics.
Those metrics are accessible through an RPC API.

Usage :
	monitord [-config path]
where path is the relative path to the config file of the daemon.
If the config flag is not provided, monitord will look for
a file named config.json in the current directory.

Note that monitord's config file is different from monitorctl's.
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
