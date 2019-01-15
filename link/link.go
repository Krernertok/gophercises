package link

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
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

	traverseTree(doc)

	return []Link{}, nil
}

func traverseTree(n *html.Node) {
	child := n.FirstChild

	for child != nil {
		traverseTree(child)
		child = child.NextSibling
	}

	if n.DataAtom == atom.A {
		handleAnchorNode(n)
	}
}

// TODO: Extract href and text and return Link
func handleAnchorNode(n *html.Node) {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			fmt.Println(attr.Val)
		}
	}
}
