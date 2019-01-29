package main

import (
	"flag"
	"fmt"
	"github.com/krernertok/gophercises/sitemap"
	"log"
	"strings"
)

// Solution plan:
// 1. Get the domain
// 		- If no "http://" or "https://" in the beginning, add to URL string
// 2. Init storage for visited URLs (e.g. set using map keys)
// 3. Init storage for remaining URLs (e.g. container/list)
// 4. Get slice of links and pop them into the queue
// 5. Get URL from queue, get links, update visited URLs, pop
//    unvisited links into queue
// 6. Repeat until queue is empty
// 7. Create XML from visited URLs

const (
	path = "urls.xml"
)

func main() {
	depth := flag.Int("depth", 0, "maximum depth to search for links")
	flag.Parse()

	domain := flag.Arg(0)
	if domain == "" {
		fmt.Println("Provide the URL for the domain you want the sitemap for.", domain)
		return
	}

	domain = addScheme(domain)
	urls, err := sitemap.GetURLs(domain, *depth)

	if err != nil {
		log.Fatal(err)
	}

	err = sitemap.WriteLinksXML(path, urls)
	if err != nil {
		log.Fatal("Error writing XML:", err)
	}
}

func addScheme(domain string) string {
	if !strings.HasPrefix(domain, "http://") &&
		!strings.HasPrefix(domain, "https://") {
		domain = "http://" + domain
	}
	return domain
}
