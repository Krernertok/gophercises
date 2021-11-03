package cli

import (
	"bufio"
	"fmt"
	"gophercises/cyoa/data"
	"os"
	"strconv"
	"strings"
	text "text/template"
)

func RunCliMode(entry string, chapters map[string]data.Chapter) {
	data.ValidEntry(entry, chapters)

	funcMap := text.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}

	tmpl := "cyoa.txt.tmpl"
	template, err := text.New(tmpl).Funcs(funcMap).ParseFiles(tmpl)
	if err != nil {
		fmt.Println("Could not parse template:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	next := entry
chapter:
	for {
		chapter := chapters[next]
		err = template.Execute(os.Stdout, chapter)
		if err != nil {
			panic(err)
		}

		numOptions := len(chapter.Options)
		if numOptions == 0 {
			break
		}

		fmt.Println("Choose option:")
		for scanner.Scan() {
			selection := strings.TrimSpace(scanner.Text())
			i, err := strconv.Atoi(selection)
			if err != nil || i < 0 || i > numOptions {
				fmt.Println("Invalid option. Try again.")
				continue
			}

			// convert i to 0-based index
			next = chapter.Options[i-1].Arc
			continue chapter
		}

		err := scanner.Err()
		if err != nil {
			panic(err)
		} else if err == nil {
			// if user inputs EOF, exit
			break
		}
	}
}
