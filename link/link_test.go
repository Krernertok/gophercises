package link

import (
	"bufio"
	"os"
	"path"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var tests map[string][]Link = map[string][]Link{
	"ex1.html": {
		Link{"/other-page", "A link to another page"},
	},
	"ex2.html": {
		Link{"https://www.twitter.com/joncalhoun", " Check me out on twitter "},
		Link{"https://github.com/gophercises", " Gophercises is on Github! "},
	},
	"ex3.html": {
		Link{"#", "Login "},
		Link{"/lost", "Lost? Need help?"},
		Link{"https://twitter.com/marcusolsson", "@marcusolsson"},
	},
	"ex4.html": {
		Link{"/dog-cat", "dog cat "},
	},
	"ex5.html": {
		Link{"/dog-cat", "dog cat hippopotamus"},
	},
}

func TestExtractLinks(t *testing.T) {
	for filename, result := range tests {
		file, err := os.Open(path.Join("testdata", filename))
		if err != nil {
			panic(err)
		}

		fileReader := bufio.NewReader(file)
		doc, err := html.Parse(fileReader)
		if err != nil {
			panic(err)
		}

		links := []Link{}
		links = ExtractLinks(doc, links)

		for i, link := range links {
			if strings.Compare(link.Href, result[i].Href) != 0 {
				t.Errorf("Expected href: %s\nReceived: %s", result[i].Href, link.Href)
			}

			if strings.Compare(link.Text, result[i].Text) != 0 {
				t.Errorf("Expected text: %s\nReceived: %s", result[i].Text, link.Text)
			}
		}
	}
}
