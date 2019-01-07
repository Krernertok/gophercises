package main

import (
	"fmt"
	"github.com/krernertok/gophercises/cyoa/story"
	"os"
)

func main() {
	var arc story.Arc
	var ok bool

	filePath := "../data/gopher.json"
	firstArc := "intro"

	story, err := story.ParseStory(filePath)

	if err != nil {
		fmt.Println("Couldn't parse the story file:", filePath)
		os.Exit(1)
	}

	arc = story[firstArc]
	fmt.Println("\n" + arc.Title + "\n")

	for {
		for _, p := range arc.Story {
			fmt.Printf("%s\n\n", p)
		}

		numOptions := len(arc.Options)
		if numOptions == 0 {
			break
		}

		fmt.Printf("\nWhat do you want to do?\n\n")

		for i, o := range arc.Options {
			fmt.Printf("%d. %s\n", i+1, o.Text)
		}

		// subtract 1 for 0-based index
		option := getOption(numOptions) - 1
		arc, ok = story[arc.Options[option].ArcName]

		if !ok {
			fmt.Println("Invalid arc name.")
			os.Exit(1)
		}
	}
}

func getOption(numOptions int) int {
	var option int

	for {
		fmt.Print("> ")
		_, err := fmt.Scanf("%d", &option)

		if err != nil {
			printInputPrompt()
			continue
		}

		if option < 1 || option > numOptions {
			printInputPrompt()
			continue
		}

		break
	}

	return option
}

func printInputPrompt() {
	fmt.Println("Please enter the number of the option you want to pick.")
}
