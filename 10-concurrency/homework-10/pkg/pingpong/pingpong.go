package pingpong

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"
)

// Config задает параметры партии в пинг-понг.
type Config struct {
	Player1Name string        // Имя первого игрока (по умолчанию "Игрок 1")
	Player2Name string        // Имя второго игрока (по умолчанию "Игрок 2")
	TargetScore int           // Количество очков для победы в партии (по умолчанию 11)
	LuckyChance float64       // Вероятность мощного удара в угол (по умолчанию 0.20 = 20%)
	Output      io.Writer     // Куда выводить ход матча (по умолчанию os.Stdout)
	Seed        int64         // Базовый сид для детерминированных тестов (0 = time.Now)
	HitDelay    time.Duration // Задержка между ударами для наглядности (0 в тестах)
}

// Result содержит итоговые результаты сыгранной партии.
type Result struct {
	Scores map[string]int // Итоговый счет: имя игрока -> количество очков
	Winner string         // Имя победителя матча
	Rounds int            // Общее количество сыгранных подач (раундов)
}

// player представляет участника матча.
type player struct {
	name        string
	luckyChance float64
	rnd         *rand.Rand
	output      io.Writer
	delay       time.Duration
	outMu       *sync.Mutex
}

func (p *player) log(format string, args ...any) {
	if p.outMu != nil {
		p.outMu.Lock()
		defer p.outMu.Unlock()
	}
	fmt.Fprintf(p.output, format, args...)
}

// isLucky определяет, повезло ли игроку нанести победный удар (20% шанс).
func (p *player) isLucky() bool {
	return p.rnd.Float64() < p.luckyChance
}

// run запускает цикл игры игрока в отдельной горутине.
func (p *player) run(table chan string, rallyDone chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range table {
		if p.delay > 0 {
			time.Sleep(p.delay)
		}

		switch msg {
		case "begin":
			// Подача начинается командой "begin" в канал.
			// Произвольный игрок, принявший "begin", выполняет подачу и пишет "ping".
			p.log("%s: ping\n", p.name)
			table <- "ping"

		case "ping", "pong":
			// Во время удара игроку может повезти (с вероятностью 20%), и он мощно гасит в угол.
			if p.isLucky() {
				p.log("%s: мощно загасил шарик в самый угол и выиграл очко!\n", p.name)
				// В канал подается сигнал "stop", означающий завершение подачи.
				table <- "stop"
				// Оповещаем судью матча о победителе розыгрыша
				rallyDone <- p.name
			} else {
				// Обычный возврат шарика: на "ping" отвечаем "pong", на "pong" — "ping"
				reply := "pong"
				if msg == "pong" {
					reply = "ping"
				}
				p.log("%s: %s\n", p.name, reply)
				table <- reply
			}

		case "stop":
			// Подача завершена (соперник выиграл очко).
			// Игрок ждет следующую команду "begin" или закрытие канала.
		}
	}
}

// PlayGame запускает и проводит полноценную партию в пинг-понг между двумя горутинами.
func PlayGame(cfg Config) Result {
	if cfg.Player1Name == "" {
		cfg.Player1Name = "Игрок 1"
	}
	if cfg.Player2Name == "" {
		cfg.Player2Name = "Игрок 2"
	}
	if cfg.TargetScore <= 0 {
		cfg.TargetScore = 11
	}
	if cfg.LuckyChance <= 0 {
		cfg.LuckyChance = 0.20
	}
	if cfg.Output == nil {
		cfg.Output = io.Discard
	}

	var outMu sync.Mutex

	outMu.Lock()
	fmt.Fprintf(cfg.Output, "=== Начало матча по пинг-понгу: %s против %s (игра до %d очков) ===\n\n",
		cfg.Player1Name, cfg.Player2Name, cfg.TargetScore)
	outMu.Unlock()

	scores := map[string]int{
		cfg.Player1Name: 0,
		cfg.Player2Name: 0,
	}

	seed1 := cfg.Seed
	seed2 := cfg.Seed + 1000
	if cfg.Seed == 0 {
		now := time.Now().UnixNano()
		seed1 = now
		seed2 = now + 1
	}

	// Каждому игроку выделяется независимый генератор псевдослучайных чисел
	p1 := &player{
		name:        cfg.Player1Name,
		luckyChance: cfg.LuckyChance,
		rnd:         rand.New(rand.NewSource(seed1)),
		output:      cfg.Output,
		delay:       cfg.HitDelay,
		outMu:       &outMu,
	}
	p2 := &player{
		name:        cfg.Player2Name,
		luckyChance: cfg.LuckyChance,
		rnd:         rand.New(rand.NewSource(seed2)),
		output:      cfg.Output,
		delay:       cfg.HitDelay,
		outMu:       &outMu,
	}

	// Канал для перебрасывания шарика между половинами стола (небуферизованный)
	table := make(chan string)

	// Канал для фиксации окончания розыгрыша (подачи)
	rallyDone := make(chan string)

	var wg sync.WaitGroup

	// Каждому игроку соответствует поток (горутина)
	wg.Add(2)
	go p1.run(table, rallyDone, &wg)
	go p2.run(table, rallyDone, &wg)

	rounds := 0
	var winner string

	// Главный цикл матча (координатор / судья)
	for {
		rounds++
		outMu.Lock()
		fmt.Fprintf(cfg.Output, "--- Подача №%d ---\n", rounds)
		outMu.Unlock()

		// Подача начинается командой "begin" в канал
		table <- "begin"

		// Ожидаем завершения подачи (сигнал о выигранном очке)
		rallyWinner := <-rallyDone
		scores[rallyWinner]++

		outMu.Lock()
		fmt.Fprintf(cfg.Output, "Счёт подачи: %s (%d) : %s (%d)\n\n",
			cfg.Player1Name, scores[cfg.Player1Name],
			cfg.Player2Name, scores[cfg.Player2Name])
		outMu.Unlock()

		// Проверка условия окончания партии
		if scores[rallyWinner] >= cfg.TargetScore {
			winner = rallyWinner
			break
		}
	}

	// Завершаем партию: закрываем канал стола, чтобы игроки вышли из цикла
	close(table)

	// Дожидаемся корректного завершения всех горутин игроков
	wg.Wait()

	outMu.Lock()
	fmt.Fprintf(cfg.Output, "=== Партия окончена! Итоговый счёт: %s (%d) : %s (%d). Победитель: %s! ===\n",
		cfg.Player1Name, scores[cfg.Player1Name],
		cfg.Player2Name, scores[cfg.Player2Name],
		winner)
	outMu.Unlock()

	return Result{
		Scores: scores,
		Winner: winner,
		Rounds: rounds,
	}
}
