package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"go-core-4/lessons/10-concurrency/homework-10/pkg/pingpong"
)

func main() {
	player1 := flag.String("p1", "Игрок 1", "Имя первого игрока")
	player2 := flag.String("p2", "Игрок 2", "Имя второго игрока")
	targetScore := flag.Int("score", 11, "Количество очков для победы в партии")
	delayMs := flag.Int("delay", 25, "Задержка между ударами в миллисекундах для наглядности")
	flag.Parse()

	cfg := pingpong.Config{
		Player1Name: *player1,
		Player2Name: *player2,
		TargetScore: *targetScore,
		LuckyChance: 0.20, // 20% шанс мощного победного удара по условию
		Output:      os.Stdout,
		Seed:        time.Now().UnixNano(),
		HitDelay:    time.Duration(*delayMs) * time.Millisecond,
	}

	result := pingpong.PlayGame(cfg)

	fmt.Printf("\nМатч успешно завершён за %d розыгрышей!\n", result.Rounds)
}
