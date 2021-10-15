package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	problemFile, err := os.Open("problems.csv")
	if err != nil {
		panic(err)
	}
	defer problemFile.Close()

	csvReader := csv.NewReader(problemFile)
	problems, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	numProblems := len(problems)
	correctAnswers := 0
	scanner := bufio.NewScanner(os.Stdin)

	for i, problem := range problems {
		prob := problem[0]
		answer := problem[1]

		fmt.Printf("Problem #%d: %s = ", i+1, prob)
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			panic(err)
		}
		input := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if input == answer {
			correctAnswers++
		}
	}

	fmt.Printf("You answered %d out of %d problems correctly.",
		correctAnswers, numProblems)
}
