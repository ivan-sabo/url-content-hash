package util

import (
	"crypto/md5"
	"fmt"
	"io"
)

type Hash struct {
	URL  string
	Hash string
	Err  error
}

func HashMD5(c <-chan Content) <-chan Hash {
	h := make(chan Hash, 5)

	go func() {
		defer close(h)

		for s := range c {
			if s.Err != nil {
				h <- Hash{
					URL: s.URL,
					Err: s.Err,
				}

				continue
			}

			hashed := md5.Sum([]byte(s.Content))

			h <- Hash{
				Hash: fmt.Sprintf("%x", hashed),
				URL:  s.URL,
			}
		}
	}()

	return h
}

func PrintHashMap(h <-chan Hash, w io.Writer) {
	for s := range h {
		if s.Err == nil {
			fmt.Fprintf(w, "%s\t%s\n", s.URL, s.Hash)
		}
	}
}
