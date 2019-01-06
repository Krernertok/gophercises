package main

import (
	"fmt"
	"github.com/krernertok/gophercises/cyoa/story"
	"html/template"
	"net/http"
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

	if path == "/" {
		arc, found = h.story[firstArc]
	} else {
		if path[0] == '/' {
			path = path[1:]
		}
		arc, found = h.story[path]
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

/*
1. Parsing JSON into maps and structs DONE
2. HTTP server / Routing DONE
	- Name of the arc should be the path
3. Templates DONE
4. Styling
*/
