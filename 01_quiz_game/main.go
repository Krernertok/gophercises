package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var filepath *string
var defaultPath = "data/quiz.csv"

func init() {
	log.SetOutput(os.Stdout)

	pathUsage := "path to the CSV file containing the questions and answers"
	filepath = flag.String("file", defaultPath, pathUsage)
}

func main() {
	var numQuestions, numCorrect int

	flag.Parse()
	questions := getQuestions(*filepath)

	fmt.Println("Answer the following questions as quickly as possible.")
	fmt.Println("Press ENTER to start.")

	var start string
	fmt.Scanln(&start)

	for _, q := range questions {
		if len(q) != 2 {
			continue
		}

		// print the question
		fmt.Println(q[0])

		var answer string
		fmt.Scanln(&answer)
		answer = strings.ToLower(strings.TrimSpace(answer))

		// check against the correct answer
		if answer == q[1] {
			numCorrect++
		}

		numQuestions++
	}

	fmt.Println("You got", numCorrect, "out of", numQuestions, "questions correct!")
}

func getQuestions(path string) [][]string {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal("Could not open file with path", path, "Error occurred:", err)
	}

	defer file.Close()

	questions := make([][]string, 0)
	csvReader := csv.NewReader(file)

	for {
		question, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("Encountered error reading CSV file:", err)
		}

		questions = append(questions, question)
	}

	return questions
}
