package shuffler

import (
	"math/rand"
	"quiz/questions"
)

func Shuffle(list []questions.Quiz) {
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
}
