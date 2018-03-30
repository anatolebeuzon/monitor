package agent

// PollAll concurrently polls all websites once.
// It blocks until all websites have been polled.
func (w *Websites) PollAll() {
	notify := make(chan bool)
	for i := range *w {
		go (*w)[i].Poll()
	}
	for i := 0; i < len(*w); i++ {
		<-notify
	}
}

func NewWebsites(URLs []string) (w Websites) {
	for _, url := range URLs {
		w = append(w, Website{URL: url})
	}
	return
}
