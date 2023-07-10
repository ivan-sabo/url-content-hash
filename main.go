package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

var threadLimit = flag.Int("parallel", 10, "Limit of parallel requests")

type content struct {
	URL     string
	Content string
	Err     error
}

type hash struct {
	URL  string
	Hash string
	Err  error
}

func main() {
	flag.Parse()
	urls := flag.Args()

	var wg sync.WaitGroup

	content := make(chan content, 5)
	fetch(urls, int(*threadLimit), content, &wg)

	hash := make(chan hash, 5)
	go hashMD5(content, hash)
	go printHashMap(hash, os.Stdout, &wg)

	wg.Wait()
}

func fetch(urls []string, threadLimit int, c chan<- content, wg *sync.WaitGroup) {
	guard := make(chan struct{}, threadLimit)

	for _, url := range urls {
		wg.Add(1)

		guard <- struct{}{}

		go func(url string) {
			fetchURL(url, c)
			<-guard
		}(url)
	}

	//close(c)
}

func fetchURL(url string, c chan<- content) {
	resp, err := http.Get(url)
	if err != nil {
		c <- content{
			Err: err,
		}
		return
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		c <- content{
			Err: err,
		}
		return
	}

	c <- content{
		Content: string(result),
		URL:     url,
	}
}

func hashMD5(c <-chan content, h chan<- hash) {
	for s := range c {
		if s.Err != nil {
			h <- hash{
				URL: s.URL,
				Err: s.Err,
			}

			continue
		}

		hashed := md5.Sum([]byte(s.Content))

		h <- hash{
			Hash: fmt.Sprintf("%x", hashed),
			URL:  s.URL,
		}
	}

	close(h)
}

func printHashMap(h <-chan hash, w io.Writer, wg *sync.WaitGroup) {
	for s := range h {
		if s.Err != nil {
			fmt.Fprintf(w, "%s\t%s", s.URL, s.Hash)
		}
		wg.Done()
	}
}
