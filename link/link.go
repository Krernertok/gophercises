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

	links := traverseTree(doc)
	return links, nil
}

func traverseTree(n *html.Node) []Link {
	var links []Link

	child := n.FirstChild
	for child != nil {
		links = append(links, traverseTree(child)...)
		child = child.NextSibling
	}

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
		}
	}

	textFragments := textFromAnchorNode(n)
	text := strings.Join(textFragments, "")

	return Link{href, text}
}

func textFromAnchorNode(n *html.Node) []string {
	var textFragments []string

	child := n.FirstChild
	for child != nil {
		if child.Type == html.TextNode {
			textFragments = append(textFragments, child.Data)
		}
		child = child.NextSibling
	}

	return textFragments
}
