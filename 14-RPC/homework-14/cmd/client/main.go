package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"

	"go-core-4/14-RPC/homework-14/pkg/messages"
)

const serverAddr = "localhost:8080"

func main() {
	fmt.Println("=== RPC Клиент GoSearch Message Service ===")

	// 1. Подключение к серверу RPC по TCP
	client, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Ошибка подключения к RPC-серверу (%s): %v", serverAddr, err)
	}
	defer client.Close()

	fmt.Println("Успешное соединение с RPC-сервером.")

	// 2. Формируем несколько сообщений для отправки (Задача №2)
	now := time.Now()
	newBatch := []messages.Message{
		{
			ID:   1,
			Time: now.Add(-2 * time.Minute),
			Text: "Инициализация поискового робота GoSearch",
		},
		{
			ID:   2,
			Time: now.Add(-1 * time.Minute),
			Text: "Завершено индексирование корпуса документов (150 страниц)",
		},
		{
			ID:   3,
			Time: now,
			Text: "Запущен RPC-интерфейс сетевого взаимодействия",
		},
	}

	fmt.Printf("\n📤 Отправка %d сообщений на сервер (метод Service.Send)...\n", len(newBatch))
	var ackCount int
	err = client.Call("Service.Send", newBatch, &ackCount)
	if err != nil {
		log.Fatalf("Ошибка при вызове Service.Send: %v", err)
	}
	fmt.Printf("✅ Сервер подтвердил добавление сообщений: принято %d шт.\n", ackCount)

	// 3. Запрашиваем все накопленные сообщения с сервера (Задача №2)
	fmt.Println("\n📥 Запрос всех накопленных сообщений с сервера (метод Service.Messages)...")
	var allMessages []messages.Message
	err = client.Call("Service.Messages", struct{}{}, &allMessages)
	if err != nil {
		log.Fatalf("Ошибка при вызове Service.Messages: %v", err)
	}

	// 4. Вывод полученных сообщений на экран
	fmt.Printf("\n📋 Всего получено сообщений от сервера: %d\n", len(allMessages))
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%-5s | %-25s | %s\n", "ID", "ВРЕМЯ ОТПРАВКИ", "ТЕКСТ СООБЩЕНИЯ")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, m := range allMessages {
		fmt.Printf("%-5d | %-25s | %s\n", m.ID, m.Time.Format("2006-01-02 15:04:05.000"), m.Text)
	}
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Работа клиента успешно завершена.")
}
