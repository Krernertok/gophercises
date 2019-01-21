package sitemap

import (
	"github.com/krernertok/gophercises/link"
	"net/http"
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
