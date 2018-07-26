package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var countCorrect = 0
var quizLength time.Duration = 30

func main() {

	problemsCsv, _ := os.Open("problems.csv")
	reader := csv.NewReader(bufio.NewReader(problemsCsv))

	countTotal := 0
	questionAnswerPairs := make(map[string]string)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		countTotal++
		questionSlice := line[:len(line)-1]
		question := strings.Join(questionSlice, "")
		answer := line[len(line)-1]
		questionAnswerPairs[question] = answer
	}

	timer := time.NewTimer(quizLength * time.Second)
	done := make(chan bool)

	go askQuestons(done, questionAnswerPairs)

	select {
	case <-done:
	case <-timer.C:
		fmt.Println("Times up!")
	}

	fmt.Printf("Final Result: %d questions correct out of %d\n", countCorrect, countTotal)
}

func askQuestons(c chan bool, qap map[string]string) {

	for question, answer := range qap {
		fmt.Println(question)
		stdinReader := bufio.NewReader(os.Stdin)
		userAnswer, _ := stdinReader.ReadString('\n')
		userAnswer = strings.TrimSpace(userAnswer)

		if userAnswer == answer {
			fmt.Println("correct!")
			countCorrect++
		} else {
			fmt.Println("incorrect, correct answer was " + answer)
		}
	}
	c <- true
}
