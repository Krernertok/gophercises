package sitemap

import (
	"container/list"
	"encoding/xml"
	"github.com/krernertok/gophercises/link"
	"log"
	"net/http"
	"net/url"
	"os"
)

type urlElem struct {
	URL string `xml:"loc"`
}

type urlset struct {
	XMLNS string    `xml:"xmlns,attr"`
	URLs  []urlElem `xml:"url"`
}

func GetURLs(domain string) ([]string, error) {
	baseURL, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}

	unvisited := list.New()
	unvisited.PushBack(baseURL.String())

	// use urls as a set
	urls := make(map[string]struct{})
	exists := struct{}{}

	for elem := unvisited.Front(); elem != nil; elem = elem.Next() {
		href := elem.Value.(string)
		url, err := baseURL.Parse(href)

		if err != nil {
			log.Println("Skipping URL:", href, "Error:", err)
			continue
		}

		urlString := url.String()
		if _, found := urls[urlString]; found || url.Hostname() != baseURL.Hostname() {
			log.Println("Skipping:", urlString)
			continue
		}

		urls[urlString] = exists
		log.Println("Visited URL:", urlString)

		links, _ := getLinks(urlString)
		for _, link := range links {
			unvisited.PushBack(link.Href)
		}
	}

	xUrls := []string{}
	for u := range urls {
		xUrls = append(xUrls, u)
	}

	return xUrls, nil
}

func getLinks(url string) ([]link.Link, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	links, err := link.ParseLinks(res.Body)

	if err != nil {
		return nil, err
	}

	res.Body.Close()
	return links, nil
}

func WriteLinksXML(path string, urls []string) error {
	sitemap := getSitemap(urls)

	xmlFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)

	if err != nil {
		return err
	}

	defer xmlFile.Close()

	_, err = xmlFile.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	err = xml.NewEncoder(xmlFile).Encode(sitemap)

	if err != nil {
		return err
	}

	return nil
}

func getSitemap(urls []string) urlset {
	sitemap := urlset{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []urlElem{},
	}

	for _, u := range urls {
		sitemap.URLs = append(sitemap.URLs, urlElem{u})
	}

	return sitemap
}
