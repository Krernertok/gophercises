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
	randomize = flag.Bool("shuffle", defaultRandomize, randomizeUsage)

	flag.Parse()
}

func main() {
	var numCorrect int
	questions := getQuestions(*filepath, *randomize)

	printStart(*limit)

	// read the first press of ENTER
	var start string
	fmt.Scanln(&start)

	answerChan := make(chan string)
	go getAnswers(answerChan)

	for _, q := range questions {
		// just skip invalid question rows
		if len(q) != 2 {
			continue
		}

		question := q[0]
		correctAnswer := q[1]

		fmt.Println(question)

		timeout := make(chan bool)
		go setTimer(*limit, timeout)

		select {
		case answer := <-answerChan:
			if answer == correctAnswer {
				numCorrect++
			}
		case <-timeout:
		}
	}

	printEnd(numCorrect, len(questions))
}

// UI part
func printStart(timelimit int) {
	fmt.Println("Answer the following questions as quickly as possible.",
		"You have", timelimit, "seconds per answer.")
	fmt.Println("Press ENTER to start.")
}

func printEnd(correct, total int) {
	fmt.Println("You got", correct, "out of", total, "questions correct!")
}

// read questions from file
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
		randomizeQuestions(&questions)
	}

	return questions
}

func randomizeQuestions(questions *[][]string) {
	rand.Seed(time.Now().Unix())
	q := *questions
	rand.Shuffle(len(q), func(i, j int) {
		q[i], q[j] = q[j], q[i]
	})
}

// using channels
func getAnswers(answers chan<- string) {
	var input string

	for {
		fmt.Scanln(&input)
		input = strings.ToLower(strings.TrimSpace(input))
		answers <- input
	}
}

func setTimer(timeLimit int, timeout chan<- bool) {
	time.Sleep(time.Duration(timeLimit) * time.Second)
	timeout <- true
	close(timeout)
}
