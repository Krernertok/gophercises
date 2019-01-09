package main

import (
	"flag"
	"fmt"
	"github.com/krernertok/gophercises/cyoa/story"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	defaultPath := "/mnt/c/dev/goworkspace/src/github.com/krernertok/gophercises/cyoa/data/gopher.json"
	path := flag.String("file", defaultPath, "path for the story JSON file")
	firstArc := flag.String("first-arc", "intro", "name of the first story arc")
	flag.Parse()

	story, err := story.ParseStory(*path)
	if err != nil {
		fmt.Println("Couldn't parse the story file:", *path)
		return
	}

	templatePath := "/mnt/c/dev/goworkspace/src/github.com/krernertok/gophercises/cyoa/templates/arc.html"
	templateName := filepath.Base(templatePath)
	template, err := template.New(templateName).ParseFiles(templatePath)
	if err != nil {
		fmt.Println("Invalid template or path for template file:", *path)
		return
	}

	handler := getStoryHandler(story, template, *firstArc)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func getStoryHandler(story story.Story, template *template.Template, arcName string) http.Handler {
	return storyHandler{story, template, arcName}
}

type storyHandler struct {
	story    story.Story
	template *template.Template
	firstArc string
}

func (h storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var arc story.Arc
	var found bool

	path := r.URL.Path
	// this should come from somewhere else, a configuration file?
	// or should serving CSS files be done somehow differently?
	baseDir := "/mnt/c/dev/goworkspace/src/github.com/krernertok/gophercises/cyoa/"

	switch {
	case path == "/":
		arc, found = h.story[h.firstArc]
	case strings.HasPrefix(path, "/css/"):
		cssPath := baseDir + strings.TrimLeft(path, "/")
		http.ServeFile(w, r, cssPath)
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
