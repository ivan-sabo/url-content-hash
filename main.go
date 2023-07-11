// This application fetches the content of URLs provided as Command Line arguments,
// calculates MD5 hash using fetched content, and prints requested URLs and calculated
// MD5 hashes.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ivan-sabo/url-content-hash/internal/util"
)

var threadLimit = flag.Int("parallel", 10, "Limit of parallel requests")

func main() {
	flag.Parse()
	urls := urlAppendPrefix(flag.Args())

	c := util.Fetch(urls, int(*threadLimit))
	h := util.HashMD5(c)
	util.PrintHashMap(h, os.Stdout)
}

func urlAppendPrefix(urls []string) []string {
	updated := make([]string, 0, len(urls))

	for _, url := range urls {
		if !strings.HasPrefix(url, "http") {
			updated = append(updated, fmt.Sprintf("http://%s", url))
			continue
		}

		updated = append(updated, url)
	}

	return updated
}
