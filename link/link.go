package link

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"strings"
)

// Link consists of an anchor tags href attribute and the text
// text within the anchor tags.
type Link struct {
	Href string
	Text string
}

// ParseLinks parses HTML and returns a slice of Links containing
// a Link instance for every anchor tag in the HTML.
func ParseLinks(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)

	if err != nil {
		return nil, err
	}

	links := getLinks(doc)
	return links, nil
}

func getLinks(n *html.Node) []Link {
	var links []Link

	// Front end recursion
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		links = append(links, getLinks(child)...)
	}

	// Process current node
	if n.DataAtom == atom.A {
		links = append(links, handleAnchorNode(n))
	}

	return links
}

func handleAnchorNode(n *html.Node) Link {
	var href string

	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}

	text := textFromNode(n)

	return Link{href, text}
}

func textFromNode(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type != html.ElementNode {
		return ""
	}

	var text string
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		text += textFromNode(child)
	}

	return strings.Join(strings.Fields(text), " ")
}
