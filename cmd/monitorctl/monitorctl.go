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
	store := client.NewStore()
	s.Init(store)
	d := client.NewDashboard(store, &config, s.UpdateUI)
	d.Show()
}
