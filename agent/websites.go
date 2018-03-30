package agent

// PollAll concurrently polls all websites once.
// It blocks until all websites have been polled.
func (websites *Websites) PollAll() {
	notify := make(chan bool)
	for i := range *websites {
		go (*websites)[i].Poll()
	}
	for i := 0; i < len(*websites); i++ {
		<-notify
	}
}

func NewWebsites(URLs []string) (websites Websites) {
	for _, url := range URLs {
		websites = append(websites, Website{URL: url})
	}
	return
}
