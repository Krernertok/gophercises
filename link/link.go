package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func extractLinks(node *html.Node, links []Link) []Link {
	if node.FirstChild != nil {
		links = extractLinks(node.FirstChild, links)
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		var href string
		var text string

		for _, attr := range node.Attr {
			if attr.Key == "href" {
				href = attr.Val
			}
		}

		text = extractText(node.FirstChild)
		link := Link{href, text}
		links = append(links, link)
	}

	if node.NextSibling != nil {
		links = extractLinks(node.NextSibling, links)
	}

	return links
}

func extractText(node *html.Node) string {
	extractedText := ""

	if node.Type == html.TextNode {
		if extractedText != "" {
			extractedText += " "
		}
		extractedText = strings.TrimSpace(node.Data)
	}

	if node.FirstChild != nil {
		extractedText += extractText(node.FirstChild)
	}

	if node.NextSibling != nil {
		extractedText += extractText(node.NextSibling)
	}

	return extractedText
}

func main() {
	filename := flag.String("file", "ex1.html", "filename to extract links from")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	links := []Link{}
	reader := bufio.NewReader(f)
	doc, err := html.Parse(reader)
	if err != nil {
		panic(err)
	}

	links = extractLinks(doc, links)

	for _, link := range links {
		fmt.Println("Href:", link.Href)
		fmt.Println("Text:", link.Text)
		fmt.Println("---")
	}
}
