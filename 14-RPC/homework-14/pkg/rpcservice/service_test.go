package rpcservice

import (
	"net"
	"net/rpc"
	"sync"
	"testing"
	"time"

	"go-core-4/14-RPC/homework-14/pkg/messages"
)

func TestService_DirectMethods(t *testing.T) {
	svc := New()

	now := time.Now()
	testMsgs := []messages.Message{
		{ID: 1, Time: now, Text: "First message"},
		{ID: 2, Time: now.Add(time.Second), Text: "Second message"},
	}

	var addedCount int
	err := svc.Send(testMsgs, &addedCount)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if addedCount != 2 {
		t.Errorf("addedCount = %d, want 2", addedCount)
	}
	if svc.Count() != 2 {
		t.Errorf("Count() = %d, want 2", svc.Count())
	}

	var fetched []messages.Message
	err = svc.Messages(struct{}{}, &fetched)
	if err != nil {
		t.Fatalf("Messages failed: %v", err)
	}
	if len(fetched) != 2 {
		t.Fatalf("len(fetched) = %d, want 2", len(fetched))
	}
	if fetched[0].Text != "First message" || fetched[1].Text != "Second message" {
		t.Errorf("unexpected content: %v", fetched)
	}
}

func TestService_RPCOverNetPipe(t *testing.T) {
	svc := New()
	server := rpc.NewServer()
	if err := server.Register(svc); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Изолированное тестирование сетевого RPC в памяти через net.Pipe()
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	go server.ServeConn(serverConn)

	client := rpc.NewClient(clientConn)
	defer client.Close()

	// 1. Отправка сообщений клиентом
	now := time.Now()
	input := []messages.Message{
		{ID: 101, Time: now, Text: "RPC Hello"},
		{ID: 102, Time: now.Add(time.Millisecond * 100), Text: "RPC World"},
	}

	var ack int
	err := client.Call("Service.Send", input, &ack)
	if err != nil {
		t.Fatalf("client.Call(Service.Send) failed: %v", err)
	}
	if ack != 2 {
		t.Errorf("ack = %d, want 2", ack)
	}

	// 2. Получение накопленных сообщений
	var result []messages.Message
	err = client.Call("Service.Messages", struct{}{}, &result)
	if err != nil {
		t.Fatalf("client.Call(Service.Messages) failed: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("len(result) = %d, want 2", len(result))
	}
	if result[0].ID != 101 || result[1].ID != 102 {
		t.Errorf("unexpected messages: %+v", result)
	}

	// 3. Асинхронный вызов через client.Go
	asyncInput := []messages.Message{
		{ID: 103, Time: time.Now(), Text: "Async message"},
	}
	var asyncAck int
	call := client.Go("Service.Send", asyncInput, &asyncAck, nil)
	<-call.Done
	if call.Error != nil {
		t.Fatalf("async call failed: %v", call.Error)
	}

	// Проверяем итоговое количество
	var finalResult []messages.Message
	_ = client.Call("Service.Messages", struct{}{}, &finalResult)
	if len(finalResult) != 3 {
		t.Errorf("len(finalResult) = %d, want 3", len(finalResult))
	}
}

func TestService_Concurrency(t *testing.T) {
	svc := New()
	server := rpc.NewServer()
	_ = server.Register(svc)

	var wg sync.WaitGroup
	workers := 20
	messagesPerWorker := 5

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			sConn, cConn := net.Pipe()
			defer sConn.Close()
			defer cConn.Close()

			go server.ServeConn(sConn)
			c := rpc.NewClient(cConn)
			defer c.Close()

			var batch []messages.Message
			for i := 0; i < messagesPerWorker; i++ {
				batch = append(batch, messages.Message{
					ID:   workerID*1000 + i,
					Time: time.Now(),
					Text: "Concurrent message",
				})
			}

			var count int
			if err := c.Call("Service.Send", batch, &count); err != nil {
				t.Errorf("Send error: %v", err)
			}

			var msgs []messages.Message
			if err := c.Call("Service.Messages", struct{}{}, &msgs); err != nil {
				t.Errorf("Messages error: %v", err)
			}
		}(w)
	}

	wg.Wait()

	wantTotal := workers * messagesPerWorker
	if svc.Count() != wantTotal {
		t.Errorf("Total count = %d, want %d", svc.Count(), wantTotal)
	}
}

// BenchmarkService_DirectSend измеряет скорость прямого добавления сообщений в сервис.
func BenchmarkService_DirectSend(b *testing.B) {
	svc := New()
	batch := []messages.Message{
		{ID: 1, Time: time.Now(), Text: "Benchmark message"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var ack int
		_ = svc.Send(batch, &ack)
	}
}

// BenchmarkService_DirectMessages измеряет скорость извлечения сообщений из сервиса.
func BenchmarkService_DirectMessages(b *testing.B) {
	svc := New()
	for i := 1; i <= 50; i++ {
		var ack int
		_ = svc.Send([]messages.Message{{ID: i, Time: time.Now(), Text: "Test"}}, &ack)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var fetched []messages.Message
		_ = svc.Messages(struct{}{}, &fetched)
	}
}

// BenchmarkService_RPCCall_NetPipe измеряет сквозную скорость RPC-вызова по каналу net.Pipe.
func BenchmarkService_RPCCall_NetPipe(b *testing.B) {
	svc := New()
	server := rpc.NewServer()
	_ = server.Register(svc)

	sConn, cConn := net.Pipe()
	defer sConn.Close()
	defer cConn.Close()

	go server.ServeConn(sConn)
	client := rpc.NewClient(cConn)
	defer client.Close()

	batch := []messages.Message{
		{ID: 1, Time: time.Now(), Text: "RPC Benchmark message"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var ack int
		_ = client.Call("Service.Send", batch, &ack)
	}
}

