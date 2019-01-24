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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide the URL for the domain you want the sitemap for.")
		return
	}

	domain := os.Args[1]
	if !strings.HasPrefix(domain, "http://") &&
		!strings.HasPrefix(domain, "https://") {
		domain = "http://" + domain
	}
	baseURL, err := url.Parse(domain)

	if err != nil {
		fmt.Println(err)
		return
	}

	visited := make(map[string]struct{})
	exists := struct{}{}
	unvisited := list.New()

	links, err := sitemap.GetLinks(baseURL.String())

	if err != nil {
		fmt.Println(err)
		return
	}

	visited[baseURL.String()] = exists

	for _, link := range links {
		unvisited.PushBack(link)
	}

	firstElement := unvisited.Front()

	for firstElement != nil {
		href := firstElement.Value.(link.Link).Href
		url, err := baseURL.Parse(href)

		if err != nil {
			fmt.Println("Skipping href:", href, "Error:", err, "URL:", url)
			unvisited.Remove(firstElement)
			firstElement = unvisited.Front()
			continue
		}

		urlString := url.String()

		if _, found := visited[urlString]; found {
			fmt.Println("Skipping:", urlString)
			unvisited.Remove(firstElement)
			firstElement = unvisited.Front()
			continue
		}

		links, _ = sitemap.GetLinks(urlString)
		visited[urlString] = exists

		fmt.Println("Visited URL:", urlString)

		for _, link := range links {
			unvisited.PushBack(link)
		}

		unvisited.Remove(firstElement)
		firstElement = unvisited.Front()
	}
}
