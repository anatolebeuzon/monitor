package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Server     string
	Statistics struct {
		Left  Statistic
		Right Statistic
	}
	Alerts struct {
		Frequency int
		Timespan  int
	}
}

type Statistic struct {
	Frequency int
	Timespan  int
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
