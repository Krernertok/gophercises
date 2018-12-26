package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var filepath *string
var limit *int

func init() {
	log.SetOutput(os.Stdout)

	defaultPath := "data/quiz.csv"
	pathUsage := "path to the CSV file containing the questions and answers"
	filepath = flag.String("file", defaultPath, pathUsage)

	defaultLimit := 30
	limitUsage := "time limit for answering a question"
	limit = flag.Int("limit", defaultLimit, limitUsage)
}

func main() {
	var numCorrect int

	flag.Parse()
	questions := getQuestions(*filepath)
	numQuestions := len(questions)

	fmt.Println("Answer the following questions as quickly as possible.",
		"You have", *limit, "seconds per answer.")
	fmt.Println("Press ENTER to start.")

	var start string
	fmt.Scanln(&start)

	for _, q := range questions {
		if len(q) != 2 {
			continue
		}

		answerChan := make(chan string)
		timeout := make(chan bool)

		go askQuestion(q[0], answerChan)
		go timer(*limit, timeout)

		select {
		case answer := <-answerChan:
			if answer == q[1] {
				numCorrect++
			}
		case <-timeout:
			fmt.Println("Time's up! You got", numCorrect, "out of", numQuestions, "questions correct.")
			os.Exit(0)
		}
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

func askQuestion(question string, answer chan<- string) {
	var input string
	fmt.Println(question)

	fmt.Scanln(&input)
	input = strings.ToLower(strings.TrimSpace(input))

	answer <- input
}

func timer(timeLimit int, timeout chan<- bool) {
	time.Sleep(time.Duration(timeLimit) * time.Second)
	timeout <- true
}
