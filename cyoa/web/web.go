package web

import (
	"fmt"
	"gophercises/cyoa/data"
	html "html/template"
	"net/http"
	"strings"
)

type ChapterHandler struct {
	chapters map[string]data.Chapter
	template *html.Template
	entry    string
}

func (c ChapterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = c.entry
	}

	chapter, found := c.chapters[path]
	if !found {
		fmt.Fprintf(w, "Chapter not found.")
	}

	err := c.template.Execute(w, chapter)
	if err != nil {
		panic(err)
	}
}

func RunWebMode(entry string, chapters map[string]data.Chapter) {
	data.ValidEntry(entry, chapters)

	template, err := html.New("cyoa.html.tmpl").ParseFiles("cyoa.html.tmpl")
	if err != nil {
		fmt.Println("Could not parse template:", err)
	}

	chapterHandler := ChapterHandler{chapters, template, entry}
	http.Handle("/", chapterHandler)

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
