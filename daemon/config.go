/*
This file is used to generate a Config object from a JSON config file.
*/

package agent

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config represents the user-defined configuration of the daemon.
type Config struct {
	ListeningPort  int        // Port on which the RPC server listens
	Poll           PollConfig // See PollConfig struct below
	AlertThreshold float64    // Availability threshold that should trigger an alert when crossed
	URLs           []string   // List of URLs of websites to poll
}

// PollConfig defines the behavior of the poller.
type PollConfig struct {
	Interval        int // Interval between two polls, for each websites
	RetainedResults int // Number of poll results that should be kept. If set to 0, no poll result is ever deleted
}

// ReadConfig reads the config file and returns the associated Config object.
//
// The program exits if an error is encountered while reading the config file.
func ReadConfig(path string) Config {
	// Get working directory to resolve relative path
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Read file content
	data, err := ioutil.ReadFile(wd + "/" + path)
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
