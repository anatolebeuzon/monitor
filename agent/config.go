package agent

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	ListeningPort int
	PollInterval  int
	URLs          []string
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

// const maxWebsites = 1

// func CSVtoWebsites(path string) (websites Websites, err error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	r := csv.NewReader(bufio.NewReader(file))
// 	count := 0
// 	for {
// 		line, err := r.Read()
// 		if err == io.EOF || count == maxWebsites {
// 			break
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		website := Website{Hostname: line[1]}
// 		websites = append(websites, website)
// 		count++
// 	}
// 	return
// }
