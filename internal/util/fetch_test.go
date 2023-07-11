package util

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	expected1 := "dummy response"
	srv1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected1)
	}))
	defer srv1.Close()

	expected2 := "dummy response 2"
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected2)
	}))
	defer srv2.Close()

	threadLimit := 3
	urls := []string{srv1.URL, srv2.URL}
	c := Fetch(urls, threadLimit)

	results := make(map[string]string)
	for r := range c {
		results[r.URL] = r.Content
	}

	if len(results) != 2 {
		t.Fatalf("expected exactly two results, got: %v", results)
	}
	if results[srv1.URL] != expected1 {
		t.Fatalf("expected: %s, got: %s", expected1, results[srv1.URL])
	}
	if results[srv2.URL] != expected2 {
		t.Fatalf("expected: %s, got: %s", expected2, results[srv2.URL])
	}
}

func TestFetchInvalidURL(t *testing.T) {
	urls := []string{"test.url"}
	threadLimit := 3
	c := Fetch(urls, threadLimit)

	results := make(map[string]error)
	for r := range c {
		results[r.URL] = r.Err
	}

	if len(results) != 1 {
		t.Fatalf("expected exactly one result, got: %v", results)
	}
	if results[urls[0]] == nil {
		t.Fatalf("expected error, got: %v", results)
	}
}
