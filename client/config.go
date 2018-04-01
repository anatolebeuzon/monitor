package client

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Server     string
	Statistics []Statistic
	Alerts     struct {
		Frequency int
		Timespan  int
	}
}

type Statistic struct {
	Frequency int
	Timespan  int
}

func readConfig(path string) (config Config, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return config, err
	}

	data, err := ioutil.ReadFile(wd + "/" + path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}
