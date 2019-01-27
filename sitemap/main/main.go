package main

import (
	"container/list"
	"fmt"
	"github.com/krernertok/gophercises/link"
	"github.com/krernertok/gophercises/sitemap"
	"net/url"
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

var exists = struct{}{}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide the URL for the domain you want the sitemap for.")
		return
	}

	domain := addScheme(os.Args[1])
	baseURL, err := url.Parse(domain)

	if err != nil {
		fmt.Println(err)
		return
	}

	links, err := sitemap.GetLinks(baseURL.String())

	if err != nil {
		fmt.Println(err)
		return
	}

	unvisited := list.New()
	addLinksToList(links, unvisited)

	urls := make(map[string]struct{})
	urls[baseURL.String()] = exists

	for elem := unvisited.Front(); elem != nil; elem = elem.Next() {
		href := elem.Value.(link.Link).Href
		url, err := baseURL.Parse(href)

		if err != nil {
			fmt.Println("Skipping href:", href, "Error:", err)
			continue
		}

		urlString := url.String()

		if _, found := urls[urlString]; found || url.Hostname() != baseURL.Hostname() {
			fmt.Println("Skipping:", urlString)
			continue
		}

		urls[urlString] = exists

		fmt.Println("Visited URL:", urlString)

		links, _ = sitemap.GetLinks(urlString)
		addLinksToList(links, unvisited)
	}

	xUrls := []string{}
	for u := range urls {
		xUrls = append(xUrls, u)
	}

	err = sitemap.WriteLinksXML(path, xUrls)
	if err != nil {
		fmt.Println("Error writing XML:", err)
	}
}

func addScheme(domain string) string {
	if !strings.HasPrefix(domain, "http://") &&
		!strings.HasPrefix(domain, "https://") {
		domain = "http://" + domain
	}
	return domain
}

func addLinksToList(links []link.Link, l *list.List) {
	for _, link := range links {
		l.PushBack(link)
	}
}
