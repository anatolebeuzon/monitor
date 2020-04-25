/*
This file is used to generate a Config object from a JSON config file.
*/

package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config represents the user-defined configuration of the daemon.
type Config struct {
	Server     string // Address on which monitord listens
	Statistics struct {
		Left  TimeConf // Left side of the dashboard
		Right TimeConf // Right side
	}
	Alerts TimeConf
}

// TimeConf defines how the client should poll the daemon
// for a specific piece of information (e.g. latest alerts).
type TimeConf struct {
	Frequency int // Frequency (in seconds) at which the daemon should be polled
	Timespan  int // Timespan (in seconds) over which metrics should be aggregated
}

// ReadConfig reads the config file and returns the associated Config object.
//
// The program exits if an error is encountered while reading the config file.
func ReadConfig(path string) Config {
	if path == "" { // Switch to default path
		path = os.Getenv("GOPATH") + "/src/github.com/anatolebeuzon/monitor/cmd/monitorctl/config.json"
	} else {
		// Get working directory to resolve relative path
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		// Create absolute path from relative path
		path = wd + "/" + path
	}

	// Read file content
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal file content in a Config object
	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
