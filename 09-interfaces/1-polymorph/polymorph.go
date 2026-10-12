// Package main демонстрирует основы полиморфизма через интерфейсы в Go,
// утиную типизацию (duck typing) и передачу абстракций (serializer, io.Writer).
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

// serializer — абстрактный интерфейсный тип данных.
// Он объявляет контракт: любой тип, имеющий метод serialize() ([]byte, error),
// неявно реализует данный интерфейс (утиная типизация: "if it quacks like a duck, it's a duck").
type serializer interface {
	serialize() ([]byte, error) // контракт интерфейса
}

// person — конкретный пользовательский тип данных (структура).
type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// serialize реализует контракт интерфейса serializer для *person (сериализация в JSON).
// Обратите внимание на pointer receiver (*person): метод привязан к указателю.
func (p *person) serialize() ([]byte, error) {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("ошибка JSON-сериализации: %w", err)
	}
	return b, nil
}

// car — другой конкретный тип данных.
type car struct {
	Model string `xml:"model"`
	Year  int    `xml:"year"`
}

// serialize реализует контракт интерфейса serializer для *car (сериализация в XML).
func (c *car) serialize() ([]byte, error) {
	b, err := xml.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("ошибка XML-сериализации: %w", err)
	}
	return b, nil
}

// store — полиморфная функция.
// Она принимает:
//  1. s serializer — любой объект, умеющий сериализовать себя в байты.
//  2. w io.Writer — любой поток вывода (файл, сетевой сокет, HTTP-ответ, буфер памяти bytes.Buffer).
//
// ВАЖНОЕ ПРАВИЛО ПРОЕКТИРОВАНИЯ В GO:
// Если функция принимает io.Writer, она НЕ должна закрывать его (Close), так как вызывающий код
// владеет ресурсом (например, если передать os.Stdout, закрытие сломает весь консольный вывод программы!).
func store(s serializer, w io.Writer) error {
	b, err := s.serialize()
	if err != nil {
		return err
	}

	// Записываем полученные байты и перевод строки
	if _, err = w.Write(b); err != nil {
		return fmt.Errorf("ошибка записи в io.Writer: %w", err)
	}
	_, err = w.Write([]byte("\n"))
	return err
}

func main() {
	p := person{
		Name: "Курт",
		Age:  27,
	}

	// Переменной интерфейсного типа можно присвоить указатель на person,
	// так как *person реализует метод serialize().
	var s serializer = &p

	// 1. Полиморфная запись в буфер памяти bytes.Buffer (вместо захламления диска)
	var buf bytes.Buffer
	if err := store(s, &buf); err != nil {
		log.Fatal(err)
	}
	fmt.Println("--- Запись person в bytes.Buffer (JSON) ---")
	fmt.Print(buf.String())

	// 2. Полиморфная запись в os.Stdout (консоль)
	fmt.Println("--- Запись person напрямую в os.Stdout ---")
	if err := store(&p, os.Stdout); err != nil {
		log.Fatal(err)
	}

	// 3. Другой тип (car), сериализуемый в XML через тот же самый интерфейс serializer
	c := car{
		Model: "Volkswagen Beetle",
		Year:  1970,
	}

	fmt.Println("--- Запись car в os.Stdout (XML) через тот же контракт ---")
	if err := store(&c, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

