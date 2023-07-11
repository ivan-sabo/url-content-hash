package util

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"testing"
)

func TestHashMD5(t *testing.T) {
	c := make(chan Content, 2)

	content1 := Content{
		URL:     "testurl1.com",
		Content: "content1",
	}
	c <- content1

	content2 := Content{
		URL:     "testurl2.com",
		Content: "content2",
	}
	c <- content2
	close(c)

	h := HashMD5(c)

	results := make(map[string]string)
	for r := range h {
		results[r.URL] = r.Hash
	}

	if len(results) != 2 {
		t.Fatalf("expected exactly two results, got: %v", results)
	}
	expected1 := fmt.Sprintf("%x", md5.Sum([]byte(content1.Content)))
	if results[content1.URL] != expected1 {
		t.Fatalf("expected: %v, got: %v", expected1, results[content1.URL])
	}
	expected2 := fmt.Sprintf("%x", md5.Sum([]byte(content2.Content)))
	if results[content2.URL] != expected2 {
		t.Fatalf("expected: %v, got: %v", expected2, results[content2.URL])
	}
}

func TestWriteHash(t *testing.T) {
	h := make(chan Hash, 2)

	hash1 := Hash{
		URL:  "url1",
		Hash: "hash1",
	}
	hash2 := Hash{
		URL:  "url2",
		Hash: "hash2",
	}
	h <- hash1
	h <- hash2
	close(h)

	buff := bytes.Buffer{}
	WriteHash(h, &buff)
	s := bufio.NewScanner(&buff)

	s.Scan()
	expected1 := fmt.Sprintf("%s\t%s", hash1.URL, hash1.Hash)
	if got := s.Text(); expected1 != got {
		t.Fatalf("expected: %v, got: %v", expected1, got)
	}

	s.Scan()
	expected2 := fmt.Sprintf("%s\t%s", hash2.URL, hash2.Hash)
	if got := s.Text(); expected2 != got {
		t.Fatalf("expected: %v, got: %v", expected2, got)
	}
}
