package main

import (
	"io"
	"log"
	"os"
)

// sendReversedString разворачивает строку и отправляет её в любой приемник, реализующий io.Writer.
// Полиморфизм io.Writer позволяет использовать один и тот же метод для:
//   - записи в файл (*os.File);
//   - вывода в консоль (os.Stdout);
//   - отправки по сети (net.Conn);
//   - записи в буфер памяти (bytes.Buffer).
func sendReversedString(s string, w io.Writer) error {
	// В Go строки закодированы в UTF-8, где кириллица занимает по 2 байта на букву.
	// Если разворачивать строку как []byte, байты каждого символа перепутаются и возникнет "кракозябра".
	// Преобразование в []rune разбивает строку на отдельные кодовые точки (Unicode code points / символы).
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// w.Write принимает байты ([]byte), поэтому преобразуем перевернутые руны обратно в []byte.
	_, err := w.Write([]byte(string(runes)))
	return err
}

func main() {
	f, err := os.Create("./strings.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// 1. Запись в файл: f (*os.File) реализует io.Writer
	err = sendReversedString("АБЫРВАЛГ\n", f)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Вывод в консоль: os.Stdout также реализует io.Writer
	err = sendReversedString("АБЫРВАЛГ\n", os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

// get демонстрирует чтение всех данных из любого источника, реализующего io.Reader.
// В Go 1.16+ стандартным способом является io.ReadAll(r).
func get(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
