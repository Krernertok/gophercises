package main

import (
	"bufio"
	"flag"
	"fmt"
	"gophercises/cyoa/data"
	html "html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	text "text/template"
)

const (
	web = "web"
	cli = "cli"
)

type Mode string

func (m Mode) Set(s string) error {
	switch s {
	case web:
		return nil
	case cli:
		return nil
	default:
		return fmt.Errorf("%s is not a valid mode (web, cli)", s)
	}
}

type ChapterHandler struct {
	chapters map[string]data.Chapter
	template *html.Template
	entry    string
}

func validEntry(entry string, chapters map[string]data.Chapter) {
	if _, found := chapters[entry]; !found {
		fmt.Printf("Arc named '%s' not found in story data.", entry)
		os.Exit(1)
	}
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

func runWebMode(entry string, chapters map[string]data.Chapter) {
	validEntry(entry, chapters)

	template, err := html.New("cyoa.html.tmpl").ParseFiles("cyoa.html.tmpl")
	if err != nil {
		fmt.Println("Could not parse template:", err)
	}

	chapterHandler := ChapterHandler{chapters, template, entry}
	http.Handle("/", chapterHandler)

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func runCliMode(entry string, chapters map[string]data.Chapter) {
	validEntry(entry, chapters)

	funcMap := text.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}

	template, err := text.New("cyoa.txt.tmpl").Funcs(funcMap).ParseFiles("cyoa.txt.tmpl")
	if err != nil {
		fmt.Println("Could not parse template:", err)
	}

	// 1. loop
	// 2. print chapter
	// 2.1. if no options, exit
	// 3. scan text
	// 4. validate -> try again
	// 5. loop again
	scanner := bufio.NewScanner(os.Stdin)
	next := entry
chapter:
	for next != "" {
		chapter := chapters[next]
		err = template.Execute(os.Stdout, chapter)
		if err != nil {
			panic(err)
		}

		numOptions := len(chapter.Options)
		if numOptions == 0 {
			os.Exit(0)
		}

		fmt.Println("Choose option:")
		for scanner.Scan() {
			selection := strings.TrimSpace(scanner.Text())
			i, err := strconv.Atoi(selection)
			if err != nil {
				fmt.Println("Invalid option. Try again.")
				continue
			}

			// convert to 0-based index
			i -= 1
			if i < 0 || i >= numOptions {
				fmt.Println("Invalid option. Try again.")
				continue
			}

			next = chapter.Options[i].Arc
			continue chapter
		}

		err := scanner.Err()
		if err != nil {
			panic(err)
		} else if err == nil {
			next = ""
		}
	}
}

func main() {
	mode := flag.String("mode", "cli", "operating mode ('web' or 'cli')")
	entry := flag.String("entry", "intro", "first arc of the story")
	flag.Parse()

	chapters := data.ChaptersFromJson("gopher.json")

	switch *mode {
	case "web":
		runWebMode(*entry, chapters)
	case "cli":
		runCliMode(*entry, chapters)
	default:
		fmt.Printf("%s is not a valid mode ('web' or default 'cli')\n", *mode)
		os.Exit(1)
	}

}
