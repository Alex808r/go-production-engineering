// Package main разбирает классическую задачу с собеседований по Go:
// «Когда переменная интерфейса равна nil?» и коварство типизированного nil (typed nil).
package main

import (
	"fmt"
)

func main() {
	// 0: Чистый nil-интерфейс (eface: type = nil, data = nil).
	// Переменная типа any объявлена, но ей не присвоено значение.
	var iface any // 0

	// 1: Переменная функционального типа со значением по умолчанию (nil).
	// При передаче в цикл упаковывается в any -> (type = func(), data = nil).
	var f func() // 1

	// 2: Неинициализированная map со значением по умолчанию (nil).
	// При передаче в цикл упаковывается в any -> (type = map[string]int, data = nil).
	var m map[string]int // 2

	// 3: Типизированный указатель со значением nil.
	var p *int = nil
	// При присваивании указателя интерфейсу:
	// interface получает конкретный ТИП (*int) и нулевое значение данных (nil).
	var iface2 any = p // 3

	fmt.Println("--- Индексы элементов, у которых (iface == nil): ---")
	cycle(iface, f, m, iface2)

	fmt.Println("\n--- Демонстрация классической ошибки с возвратом error ---")
	demoTypedNilError()
}

// cycle проверяет каждый переданный аргумент на равенство nil.
//
// ВАЖНЕЙШЕЕ ПРАВИЛО GO:
// Значение интерфейса равно nil ТОЛЬКО ТОГДА, когда ОБА его компонента равны nil:
//   - указатель на метаданные типа (_type / itab) == nil
//   - указатель на сами данные (data) == nil
//
// Если переменная имела конкретный тип (func(), map, *int), то при передаче в `any`
// интерфейс запоминает этот тип, поэтому условие `iface == nil` становится ЛОЖНЫМ!
func cycle(ifaces ...any) {
	for i, iface := range ifaces {
		if iface == nil {
			fmt.Printf("Индекс %d равен nil (чистый nil-интерфейс)\n", i)
		} else {
			fmt.Printf("Индекс %d НЕ равен nil! Динамический тип: %T, значение: %v\n", i, iface, iface)
		}
	}
}

// CustomError — пользовательский тип ошибки.
type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

// getErrorBuggy — функция с классическим антипаттерном.
// Возвращает интерфейс error, но внутри создает переменную типизированного указателя *CustomError.
func getErrorBuggy(hasError bool) error {
	var err *CustomError = nil
	if hasError {
		err = &CustomError{Message: "критическая ошибка"}
	}
	// ОШИБКА: возвращается интерфейс с типом (*CustomError) и data = nil.
	// Для вызывающего кода `err != nil` ВСЕГДА будет true!
	return err
}

// getErrorCorrect — корректная реализация.
// Если ошибки нет, возвращается явный нетипизированный nil.
func getErrorCorrect(hasError bool) error {
	if hasError {
		return &CustomError{Message: "критическая ошибка"}
	}
	return nil // явно возвращаем nil-интерфейс
}

func demoTypedNilError() {
	err1 := getErrorBuggy(false)
	if err1 != nil {
		fmt.Printf("Ловушка! getErrorBuggy(false) != nil! Тип ошибки: %T, значение: %v\n", err1, err1)
	}

	err2 := getErrorCorrect(false)
	if err2 == nil {
		fmt.Println("Успех: getErrorCorrect(false) корректно вернул чистый nil (ошибки нет).")
	}
}

