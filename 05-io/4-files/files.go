package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// ВАЖНОЕ ЗАМЕЧАНИЕ О ДЕСКРИПТОРАХ:
	// os.Create открывает файл и возвращает дескриптор *os.File.
	// Вызов ниже os.WriteFile заново открывает файл по имени.
	// В реальном коде лучше либо писать прямо в созданный файл: f.Write([]byte("Текст")),
	// либо сразу вызывать os.WriteFile без предварительного os.Create.
	f, err := os.Create("./file.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	err = os.WriteFile(f.Name(), []byte("Текст"), 0666)
	if err != nil {
		log.Println(err)
		return
	}

	// os.ReadFile читает весь файл целиком в один срез байт за один вызов (удобно для небольших файлов).
	data, err := os.ReadFile(f.Name())
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Данные файла (os.ReadFile):\n%s\n", data)

	// Потоковое чтение с помощью интерфейса io.Reader (*os.File реализует io.Reader)
	var file *os.File
	file, err = os.Open(f.Name())
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	// ЧТО ЗДЕСЬ НЕ ПРАВИЛЬНО?
	// 1. Недостаточный размер буфера: в кодировке UTF-8 русское слово "Текст" занимает 10 байт
	//    (по 2 байта на букву). Буфер размером 6 байт прочитает только 3 буквы ("Тек")!
	// 2. Отсутствие цикла чтения: метод file.Read() НЕ обязан читать файл целиком за один вызов.
	//    Он читает порцию данных до размера буфера (len(buf)).
	//    Чтобы прочитать всё содержимое потока, нужно либо вызывать Read в цикле до io.EOF,
	//    либо использовать io.ReadAll(file).
	buf := make([]byte, 6)
	n, err := file.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}
	buf = buf[:n]
	fmt.Printf("Данные файла (неполные, прочитано только %d из 10 байт):\n%s\n", n, buf)
}
