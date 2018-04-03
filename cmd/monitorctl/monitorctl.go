/*
monitorctl is a client for the monitord daemon.
It displays website statistics and alerts on a console dashboard.

Usage :
	monitorctl [-config path]
where path is the relative path to the config file of the client.
If the config flag is not provided, monitorctl will look for
a file named config.json in the current directory.

Note that monitorctl's config file is different from monitord's.
*/
package main

import (
	"flag"
	"log"
	"monitor/client"
)

func main() {
	// Load config
	configPath := flag.String("config", "config.json", "Config file in JSON format")
	flag.Parse()
	config, err := client.ReadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	s := client.NewScheduler(config)
	store := client.NewStore()
	s.Init(store)
	d := client.NewDashboard(store, &config, s.UpdateUI)
	d.Show()
}
