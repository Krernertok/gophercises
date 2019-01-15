package main

import (
	"fmt"
	"github.com/krernertok/gophercises/link"
	"os"
)

func main() {
	path := "data/ex3.html"
	html, err := os.Open(path)

	if err != nil {
		fmt.Println("Couldn't open the HTML file:", path)
		return
	}

	link.ParseLinks(html)
}
