package util

import (
	"io"
	"net/http"
	"sync"
)

// Content holds result of the http get request to a specified URL
type Content struct {
	URL     string
	Content string
	Err     error
}

// Fetch makes concurrent http get requests to provided urls and transmits results to returning Content channel
// It also limits thread count accorting threadLimit parameter
func Fetch(urls []string, threadLimit int) <-chan Content {
	c := make(chan Content, 5)
	guard := make(chan struct{}, threadLimit)

	go func() {
		var wg sync.WaitGroup

		for _, url := range urls {
			wg.Add(1)

			guard <- struct{}{}

			go func(url string) {
				fetchURL(url, c)
				<-guard
				wg.Done()
			}(url)
		}

		go func() {
			wg.Wait()
			close(c)
		}()
	}()

	return c
}

func fetchURL(url string, c chan<- Content) {
	resp, err := http.Get(url)
	if err != nil {
		c <- Content{
			URL: url,
			Err: err,
		}
		return
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		c <- Content{
			URL: url,
			Err: err,
		}
		return
	}

	c <- Content{
		URL:     url,
		Content: string(result),
	}
}
