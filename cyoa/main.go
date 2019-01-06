package main

import (
	"fmt"
	"github.com/krernertok/gophercises/cyoa/story"
	"html/template"
	"net/http"
	"strings"
)

func main() {
	storyPath := "data/gopher.json"
	story, err := story.ParseStory(storyPath)

	if err != nil {
		panic(err)
	}

	template, err := template.New("arc.html").ParseFiles("templates/arc.html")

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", storyHandler{story, template})
}

type storyHandler struct {
	story    map[string]story.Arc
	template *template.Template
}

func (h storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var arc story.Arc
	var found bool

	firstArc := "intro"
	path := r.URL.Path

	switch {
	case path == "/":
		arc, found = h.story[firstArc]
	case strings.HasPrefix(path, "/css/"):
		http.ServeFile(w, r, strings.TrimLeft(path, "/"))
		return
	default:
		arc, found = h.story[path[1:]]
	}

	if !found {
		http.NotFound(w, r)
		return
	}

	err := h.template.Execute(w, arc)
	if err != nil {
		panic(err)
	}
}
