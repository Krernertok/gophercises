package link

import (
	"regexp"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ExtractLinks(node *html.Node, links []Link) []Link {
	if node.FirstChild != nil {
		links = ExtractLinks(node.FirstChild, links)
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		var href string
		var text string

		for _, attr := range node.Attr {
			if attr.Key == "href" {
				href = attr.Val
			}
		}

		re := regexp.MustCompile(`\s+`)
		text = extractText(node.FirstChild)
		trimmedText := re.ReplaceAllString(text, " ")

		link := Link{href, trimmedText}
		links = append(links, link)
	}

	if node.NextSibling != nil {
		links = ExtractLinks(node.NextSibling, links)
	}

	return links
}

func extractText(node *html.Node) string {
	extractedText := ""

	if node.Type == html.TextNode {
		extractedText = node.Data
	}

	if node.FirstChild != nil {
		extractedText += extractText(node.FirstChild)
	}

	if node.NextSibling != nil {
		extractedText += extractText(node.NextSibling)
	}

	return extractedText
}
