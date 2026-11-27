package game

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"quiz/questions"
	"strings"
)

func askQuestion(q questions.Quiz, reader *bufio.Reader) (bool, error) {
	fmt.Printf("Enter the Capital of %s", q.Country)

	answer, err := getInputUser(reader)
	if err != nil {
		return false, err
	}
	if answer == strings.ToLower(q.Capital) {
		fmt.Println("Correct!")
		return true, nil
	} else {
		fmt.Println("Incorrect! The correct answer is ", q.Capital)
		return false, nil
	}

}

func getInputUser(reader *bufio.Reader) (string, error) {
	for {
		fmt.Print(" Enter your answer: ")
		input, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Some error occurred while reading the input.", err)
			continue
		}

		input = strings.TrimSpace(input)

		if input != "" {
			return strings.ToLower(input), nil
		}
		if err == io.EOF {
			return "", err
		}
	}
}

func StartGame(questions []questions.Quiz) (correctAnswers, wrongAnswersCount uint, wrongAnswersList []questions.Quiz) {
	fmt.Println("Welcome to Go Quiz!")
	reader := bufio.NewReader(os.Stdin)

	for _, q := range questions {
		correct, err := askQuestion(q, reader)
		if err != nil {
			break
		}
		if correct {
			correctAnswers++
		} else {
			wrongAnswersCount++
			wrongAnswersList = append(wrongAnswersList, q)
		}
	}

	return correctAnswers, wrongAnswersCount, wrongAnswersList
}
