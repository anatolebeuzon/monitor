package main

import (
	"fmt"
	"log"
	"monitor/daemon"
	"os"
	"os/signal"
)

const (
	name        = "monitord"
	description = "Daemon that polls websites and gather related metrics"
	configPath  = "config.json"
)

func main() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive the signal.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Load config
	config, err := agent.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	websites := agent.NewWebsites(config.URLs)
	websites.SchedulePolls(config.Poll)

	// Create RPC handler
	h := &agent.Handler{
		Websites:       &websites,
		AlertThreshold: config.AlertThreshold,
		Done:           make(chan bool),
	}
	go agent.ServeRPC(h, config.ListeningPort)

	// Handle interrupt by system signal
	<-interrupt
	// TODO: properly close RPC
	fmt.Println("Closing properly...")
	return
}
