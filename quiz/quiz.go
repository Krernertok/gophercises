package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	file := flag.String("file", "problems.csv", "file with problems")
	limit := flag.Int("limit", 30, "time limit in seconds (default 30)")
	flag.Parse()

	problemFile, err := os.Open(*file)
	if err != nil {
		fmt.Println("Issue opening the problems file:", *file)
		os.Exit(1)
	}
	defer problemFile.Close()

	csvReader := csv.NewReader(problemFile)
	problems, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Could not parse the CSV file:", file)
		os.Exit(1)
	}

	numProblems := len(problems)
	correctAnswers := 0
	scanner := bufio.NewScanner(os.Stdin)
	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	answers := make(chan string)

Loop:
	for i, problem := range problems {
		prob := problem[0]
		answer := problem[1]

		fmt.Printf("Problem #%d: %s = ", i+1, prob)

		go func() {
			scanner.Scan()

			if err := scanner.Err(); err != nil {
				panic(err)
			}
			answers <- strings.TrimSpace(strings.ToLower(scanner.Text()))
		}()

		select {
		case <-timer.C:
			fmt.Println("\nYou ran out of time!")
			break Loop
		case input := <-answers:
			if input == answer {
				correctAnswers++
			}
		}
	}

	fmt.Printf("You answered %d out of %d problems correctly.",
		correctAnswers, numProblems)
}
