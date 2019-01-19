package main

import (
	"fmt"
	"github.com/krernertok/gophercises/link"
	"os"
)

func main() {
	path := "data/ex2.html"
	html, err := os.Open(path)

	if err != nil {
		fmt.Println("Couldn't open the HTML file:", path)
		return
	}

	links, err := link.ParseLinks(html)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, l := range links {
		fmt.Printf("Link: '%s', Text: '%s'\n", l.Href, l.Text)
	}
}
