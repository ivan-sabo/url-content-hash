package util

import (
	"crypto/md5"
	"fmt"
	"io"
)

// Hash holds a result of hash function of content of the specified URL
type Hash struct {
	URL  string
	Hash string
	Err  error
}

// HashMD5 calculates MD5 hash of the content of received Content object and pushes the result to
// returning Hash channel
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

// WriteHash writes the received data of Hash object to a provided writer
func WriteHash(h <-chan Hash, w io.Writer) {
	for s := range h {
		if s.Err == nil {
			fmt.Fprintf(w, "%s\t%s\n", s.URL, s.Hash)
		}
	}
}
