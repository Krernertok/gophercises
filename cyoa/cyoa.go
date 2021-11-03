package main

import (
	"flag"
	"fmt"
	"gophercises/cyoa/cli"
	"gophercises/cyoa/data"
	"gophercises/cyoa/web"
	"os"
)

const (
	webOption = "web"
	cliOption = "cli"
)

type Mode string

func (m Mode) Set(s string) error {
	switch s {
	case webOption:
		return nil
	case cliOption:
		return nil
	default:
		return fmt.Errorf("%s is not a valid mode (web, cli)", s)
	}
}

func main() {
	modeHelp := fmt.Sprintf("operating mode ('%s' or '%s'",
		webOption, cliOption)
	mode := flag.String("mode", cliOption, modeHelp)
	entry := flag.String("entry", "intro", "first arc of the story")
	flag.Parse()

	chapters := data.ChaptersFromJson("gopher.json")

	switch *mode {
	case webOption:
		web.RunWebMode(*entry, chapters)
	case cliOption:
		cli.RunCliMode(*entry, chapters)
	default:
		fmt.Printf("%s is not a valid mode ('web' or default 'cli')\n", *mode)
		os.Exit(1)
	}
}
