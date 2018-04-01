package main

import (
	"log"
	"monitor/client"
)

const configPath = "config.json"

func main() {
	// Load config
	config, err := client.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	s := client.NewScheduler(config)
	agg := s.Init()
	d := client.NewDashboard(agg, s.UpdateUI)
	d.Show()
}
