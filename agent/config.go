package agent

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	ListeningPort int
	Poll          PollConfig
	URLs          []string
}

type PollConfig struct {
	Interval        int
	RetainedResults int
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
