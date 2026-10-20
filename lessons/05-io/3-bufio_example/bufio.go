package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("./bufio.go")
	if err != nil {
		log.Fatal(err)
	}

	b, err := get(f)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()

	f, err = os.Create("./bufio_copy.go")
	if err != nil {
		log.Fatal(err)
	}

	err = store(f, b)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()
}

// store записывает байты в любой приемник io.Writer.
func store(w io.Writer, b []byte) error {
	_, err := w.Write(b)
	return err
}

// get читает данные из io.Reader построчно с помощью bufio.Scanner.
// bufio буферизует чтение в памяти (внутренний буфер по умолчанию 4 КБ),
// что значительно снижает количество системных вызовов read() к ядру ОС.
func get(r io.Reader) ([]byte, error) {
	// bufio.NewScanner разбивает входной поток на токены (по умолчанию — по строкам '\n').
	scanner := bufio.NewScanner(r)

	var b []byte
	// scanner.Scan() считывает очередную строку и возвращает true.
	// Когда поток заканчивается (EOF), Scan() возвращает false и цикл завершается.
	for scanner.Scan() {
		// scanner.Text() возвращает строку БЕЗ символа перевода строки '\n',
		// поэтому при необходимости мы добавляем "\n" обратно.
		b = append(b, []byte(scanner.Text()+"\n")...)
	}

	// Обязательная проверка: если сканирование прервалось из-за ошибки ввода-вывода,
	// scanner.Err() вернет эту ошибку (io.EOF ошибкой не считается и возвращает nil).
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return b, nil
}
