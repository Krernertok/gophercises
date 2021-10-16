package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func parseArgs() (string, int, bool) {
	file := flag.String("file", "problems.csv", "file with problems")
	limit := flag.Int("limit", 30, "time limit in seconds (default 30)")
	shuffle := flag.Bool("shuffle", false, "shuffle the problems")
	flag.Parse()
	return *file, *limit, *shuffle
}

func readCsv(file string) ([][]string, bool) {
	problemFile, err := os.Open(file)
	if err != nil {
		return nil, false
	}
	defer problemFile.Close()

	csvReader := csv.NewReader(problemFile)
	problems, err := csvReader.ReadAll()
	if err != nil {
		return nil, false
	}

	return problems, true
}

func shuffleProblems(problems [][]string) {
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
}

func main() {
	file, limit, shuffle := parseArgs()
	problems, ok := readCsv(file)
	if !ok {
		fmt.Println("Problem reading the CSV file:", file)
		os.Exit(1)
	}
	numProblems := len(problems)

	if shuffle {
		shuffleProblems(problems)
	}

	correctAnswers := 0
	scanner := bufio.NewScanner(os.Stdin)
	answers := make(chan string)

	fmt.Println("Press ENTER to start the quiz.")
	scanner.Scan()

	timer := time.NewTimer(time.Duration(limit) * time.Second)

Loop:
	for i, problem := range problems {
		prob := problem[0]
		answer := strings.ToLower(problem[1])

		fmt.Printf("Problem #%2d: %s = ", i+1, prob)

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
