/*
This file is used to generate a Config object from a JSON config file.
*/

package daemon

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config represents the user-defined configuration of the daemon.
type Config struct {
	ListeningPort int // Port on which the RPC server listens
	Default       struct {
		Interval        int     // Interval, in seconds, between two polls to a given website
		RetainedResults int     // Number of poll results that should be kept. If set to 0, no poll result is ever deleted
		Threshold       float64 // Availability threshold that should trigger an alert when crossed
	}
	Websites []WebsiteConfig // List of websites to poll
}

// WebsiteConfig represents the configuration of a specific website.
type WebsiteConfig struct {
	URL string

	// If Interval, RetainedResults or Threshold are not filled, Config.Default will be used instead
	Interval        int
	RetainedResults int
	Threshold       float64
}

// ReadConfig reads the config file and returns the associated Config object.
//
// The program exits if an error is encountered while reading the config file.
func ReadConfig(path string) Config {
	if path == "" { // Switch to default path
		path = os.Getenv("GOPATH") + "/src/github.com/oxlay/monitor/cmd/monitord/config.json"
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
