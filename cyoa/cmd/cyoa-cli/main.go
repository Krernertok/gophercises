package main

import (
	"flag"
	"fmt"
	"github.com/krernertok/gophercises/cyoa/story"
)

func main() {
	defaultPath := "/mnt/c/dev/goworkspace/src/github.com/krernertok/gophercises/cyoa/data/gopher.json"
	filepath := flag.String("file", defaultPath, "path for the story JSON file")
	firstArc := flag.String("first-arc", "intro", "name of the first story arc")
	flag.Parse()

	var arc story.Arc
	var ok bool

	story, err := story.ParseStory(*filepath)
	if err != nil {
		fmt.Println("Couldn't parse the story file:", *filepath)
		return
	}

	arc, ok = story[*firstArc]
	if !ok {
		fmt.Printf("first-arc '%s' does not match any arc in the story.\n", *firstArc)
		return
	}

	fmt.Println("\n" + arc.Title + "\n")

	for {
		fmt.Println()
		printParagraphs(arc.Paragraphs)

		// an ending arc to the story has no options, so break the loop
		if len(arc.Options) == 0 {
			break
		}

		option := getOption(arc.Options)
		arc, ok = story[option]

		if !ok {
			fmt.Println("Invalid arc name. Make sure that the story file is valid.")
			return
		}
	}
}

func getOption(options []story.Option) string {
	fmt.Printf("\nWhat do you want to do?\n\n")
	printOptions(options)

	// subtract 1 for 0-based index
	index := getInput(len(options)) - 1
	return options[index].ArcName
}

func printParagraphs(paragraphs []string) {
	for _, p := range paragraphs {
		fmt.Printf("%s\n\n", p)
	}
}

func printOptions(options []story.Option) {
	for i, o := range options {
		fmt.Printf("%d. %s\n", i+1, o.Text)
	}
}

func getInput(numOptions int) int {
	var option int

	for {
		fmt.Print("> ")
		_, err := fmt.Scan(&option)

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
