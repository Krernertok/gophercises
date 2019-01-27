package sitemap

import (
	"encoding/xml"
	"github.com/krernertok/gophercises/link"
	"net/http"
	"os"
)

func GetLinks(url string) ([]link.Link, error) {
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

type urlElem struct {
	URL string `xml:"loc"`
}

type urlset struct {
	XMLNS string    `xml:"xmlns,attr"`
	URLs  []urlElem `xml:"url"`
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
