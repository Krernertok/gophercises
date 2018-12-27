package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var filepath *string
var limit *int
var randomize *bool

func init() {
	log.SetOutput(os.Stdout)

	defaultPath := "data/quiz.csv"
	pathUsage := "path to the CSV file containing the questions and answers"
	filepath = flag.String("file", defaultPath, pathUsage)

	defaultLimit := 30
	limitUsage := "time limit for answering a question"
	limit = flag.Int("limit", defaultLimit, limitUsage)

	defaultRandomize := false
	randomizeUsage := "randomize the order of the questions"
	randomize = flag.Bool("random", defaultRandomize, randomizeUsage)
}

func main() {
	var numCorrect int

	flag.Parse()
	questions := getQuestions(*filepath, *randomize)
	numQuestions := len(questions)

	fmt.Println("Answer the following questions as quickly as possible.",
		"You have", *limit, "seconds per answer.")
	fmt.Println("Press ENTER to start.")

	var start string
	fmt.Scanln(&start)

	answerChan := make(chan string)
	go getAnswer(answerChan)

	for _, q := range questions {
		if len(q) != 2 {
			continue
		}
		fmt.Println(q[0])

		timeout := make(chan bool)
		go timer(*limit, timeout)

		select {
		case answer := <-answerChan:
			if answer == q[1] {
				numCorrect++
			}
		case <-timeout:
		}
	}

	fmt.Println("You got", numCorrect, "out of", numQuestions, "questions correct!")
}

func getQuestions(path string, randomize bool) [][]string {
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

	if randomize {
		rand.Seed(time.Now().Unix())
		rand.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}

	return questions
}

func getAnswer(answer chan<- string) {
	var input string

	for {
		fmt.Scanln(&input)
		input = strings.ToLower(strings.TrimSpace(input))
		answer <- input
	}
}

func timer(timeLimit int, timeout chan<- bool) {
	time.Sleep(time.Duration(timeLimit) * time.Second)
	timeout <- true
	close(timeout)
}
