package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

var filepath = "data/questions.csv"

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	var numQuestions, numCorrect int
	questions := getQuestions(filepath)

	for _, q := range questions {
		if len(q) != 2 {
			continue
		}
		numQuestions++

		fmt.Println(q[0])

		var answer string
		fmt.Scanln(&answer)

		if answer == q[1] {
			numCorrect++
		}
	}

	fmt.Println("You got", numCorrect, "out of", numQuestions, "correct!")
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
