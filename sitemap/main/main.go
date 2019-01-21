package main

import (
	"fmt"
	"github.com/krernertok/gophercises/sitemap"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide the URL for the domain you want the sitemap for.")
		return
	}

	links, err := sitemap.GetLinks(os.Args[1])

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(links)
}
