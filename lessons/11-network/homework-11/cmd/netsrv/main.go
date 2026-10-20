package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-core-4/lessons/07-testing/homework-07/pkg/crawler/membot"
	"go-core-4/lessons/07-testing/homework-07/pkg/engine"
	"go-core-4/lessons/07-testing/homework-07/pkg/index"
	"go-core-4/lessons/11-network/homework-11/pkg/netsrv"
)

func main() {
	fmt.Println("=== Запуск поисковой службы GoSearch (netsrv) ===")

	// 1. Инициализируем поисковый движок и робота (membot с демонстрационными данными)
	bot := membot.New()
	idx := index.New()
	eng := engine.New(bot, idx)

	// Сканируем корпус документов для индексации
	fmt.Println("Индексация демонстрационного корпуса документов...")
	if err := eng.Scan([]string{"https://go.dev"}, 1); err != nil {
		log.Fatalf("Ошибка индексации: %v", err)
	}

	docs := eng.Documents()
	fmt.Printf("Успешно проиндексировано документов: %d\n", len(docs))

	// 2. Создаем сетевую службу на порту 8000 (на всех интерфейсах 0.0.0.0:8000)
	// (Все устаревшие артефакты: флаги -s, интерактивный ввод с os.Stdin удалены)
	srv := netsrv.New(":8000", eng)

	// Перехватываем сигналы ОС для Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nПолучен сигнал завершения. Остановка сервера...")
		_ = srv.Stop()
	}()

	fmt.Println("Сетевая служба netsrv ожидает входящих подключений на :8000 (0.0.0.0:8000)...")
	fmt.Println("Для подключения используйте:")
	fmt.Println("  telnet localhost 8000")
	fmt.Println("  go run ./11-network/homework-11/cmd/netclient/main.go")

	if err := srv.Start(); err != nil {
		log.Printf("Сервер остановлен: %v", err)
	}
}
