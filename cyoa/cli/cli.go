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

var scanner = bufio.NewScanner(os.Stdin)

func scanOption(numOptions int) (int, error) {
	fmt.Println("Choose option:")
	for scanner.Scan() {
		selection := strings.TrimSpace(scanner.Text())
		i, err := strconv.Atoi(selection)
		if err != nil || i <= 0 || i > numOptions {
			fmt.Println("Invalid option. Try again.")
			continue
		}

		// convert i to 0-based index
		return i - 1, nil
	}

	err := scanner.Err()
	if err != nil {
		return 0, err
	}

	// if user inputs EOF (err == nil), return -1
	return -1, nil
}

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

	next := entry
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

		selection, err := scanOption(numOptions)
		if err != nil {
			panic(err)
		}

		if selection == -1 {
			break
		}

		next = chapter.Options[selection].Arc
	}
}
