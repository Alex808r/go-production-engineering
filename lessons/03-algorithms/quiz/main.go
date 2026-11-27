package main

import (
	"fmt"
	"os"
	"quiz/game"
	"quiz/questions"
	"quiz/shuffler"
)

func main() {
	questions, err := questions.LoadQuestions()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load questions: %v\n", err)
		os.Exit(1)
	}

	shuffler.Shuffle(questions)
	correctAnswers, wrongAnswersCount, wrongAnswersList := game.StartGame(questions)
	totalQuestions := len(questions)
	fmt.Printf("You got %d correct answers out of %d and %d wrong answers\n", correctAnswers, totalQuestions, wrongAnswersCount)
	fmt.Printf("Wrong answers: %v\n", wrongAnswersList)
}
