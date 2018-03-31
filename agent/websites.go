package agent

type Websites []Website

func NewWebsites(URLs []string) (w Websites) {
	for _, url := range URLs {
		w = append(w, Website{URL: url})
	}
	return
}

func (w Websites) schedulePolls(p PollConfig) {
	for i := range w {
		go w[i].schedulePolls(p)
	}
}
