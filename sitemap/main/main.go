package main

import (
	"fmt"
	"github.com/krernertok/gophercises/sitemap"
	"log"
	"os"
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
	if len(os.Args) < 2 {
		fmt.Println("Provide the URL for the domain you want the sitemap for.")
		return
	}

	domain := addScheme(os.Args[1])
	urls, err := sitemap.GetURLs(domain)

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
