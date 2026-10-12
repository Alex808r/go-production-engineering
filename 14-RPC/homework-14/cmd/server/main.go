package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"go-core-4/14-RPC/homework-14/pkg/rpcservice"
)

const defaultPort = ":8080"

func main() {
	fmt.Println("=== Запуск RPC-сервера накопления сообщений (Занятие №14) ===")

	// 1. Создаем сервис сообщений
	svc := rpcservice.New()

	// 2. Регистрируем сервис в RPC
	err := rpc.Register(svc)
	if err != nil {
		log.Fatalf("Ошибка регистрации RPC-сервиса: %v", err)
	}

	// 3. Открываем TCP-слушатель
	listener, err := net.Listen("tcp", defaultPort)
	if err != nil {
		log.Fatalf("Ошибка запуска TCP listener на %s: %v", defaultPort, err)
	}
	defer listener.Close()

	fmt.Printf("RPC-сервер успешно запущен и ожидает подключения на порту %s\n", defaultPort)
	fmt.Println("Поддерживаемые методы: Service.Send, Service.Messages")

	// 4. Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nПолучен сигнал завершения. Остановка RPC-сервера...")
		_ = listener.Close()
		os.Exit(0)
	}()

	// 5. Цикл обработки клиентских подключений
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-sigChan:
				return
			default:
				log.Printf("Ошибка при принятии соединения: %v", err)
				continue
			}
		}

		// Каждое клиентское подключение обслуживается в отдельной горутине
		go rpc.ServeConn(conn)
	}
}
