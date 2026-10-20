package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Пример сравнивает три способа чтения пользовательского ввода из терминала (os.Stdin):
// 1. bufio.NewReader — низкоуровневое чтение с контролем символа-разделителя.
// 2. bufio.NewScanner — идиоматичное построчное чтение строк с автоматическим удалением перевода строки.
// 3. fmt.Scanln — чтение форматированных значений (читает только до первого пробела).

func main() {
	// --- СПОСОБ 1: bufio.NewReader ---
	// Читает данные до указанного байта-разделителя ('\n').
	// Символ перевода строки остается в конце прочитанной строки, поэтому его нужно удалять вручную.
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("NewReader-> ")
		text, _ := reader.ReadString('\n')    // чтение строки (до символа перевода)
		text = strings.TrimSuffix(text, "\n") // удаление перевода строки
		text = strings.TrimSuffix(text, "\r") // удаление перевода строки для Windows CRLF
		if text == "exit" {
			break
		}
		fmt.Println("echo:", text)
	}

	// --- СПОСОБ 2: bufio.NewScanner (Рекомендуемый для построчного ввода) ---
	// Сканер автоматически разделяет поток по строкам и отсекает \r и \n.
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("NewScanner-> ")
	for scanner.Scan() {
		if scanner.Text() == "exit" {
			break
		}
		fmt.Println(scanner.Text())
		fmt.Print("NewScanner-> ")
	}

	// --- СПОСОБ 3: fmt.Scanln ---
	// Читает токены, разделенные пробелами.
	// Внимание: если пользователь введет "Hello World", в переменную s попадет только "Hello",
	// а "World" останется ждать в буфере ввода.
	for {
		fmt.Print("fmt.Scanln-> ")
		var s string
		fmt.Scanln(&s)
		if s == "exit" {
			break
		}
		fmt.Println(s)
	}
}
