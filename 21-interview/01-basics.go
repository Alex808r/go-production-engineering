package interview

import (
	"fmt"
	"sync"
)

// TypicalErrorsNoSync демонстрирует две классические проблемы:
// 1. Отсутствие синхронизации (горутины запускаются, но main завершается до их выполнения).
// 2. До Go 1.22 переменная цикла `i` захватывалась по ссылке (все горутины печатали последнее значение, обычно 10).
func TypicalErrorsNoSync() {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
}

// TypicalErrorsIndex решает проблему завершения с помощью sync.WaitGroup.
// Идиоматично передавать значение итератора как аргумент горутины: `go func(n int) { ... }(i)`.
func TypicalErrorsIndex() {
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(n int) {
			defer wg.Done()
			fmt.Println(n)
		}(i)
	}

	wg.Wait()
}

// Slices демонстрирует механику срезов в Go:
//  1. Срез s2 изначально указывает на тот же базовый массив, что и s1.
//     Поэтому s2[1] = 999 изменяет и s1[2].
//  2. При append, если превышена вместимость (cap) среза, Go выделяет НОВЫЙ базовый массив.
//     После реалокации s1 и s2 ссылаются на разные участки памяти.
func Slices() {
	s1 := []int{1, 2, 3, 4}
	s2 := s1[1:3] // [2,3]
	fmt.Println("До изменения S2")
	fmt.Println(s1) // [1 2 3 4]
	fmt.Println(s2) // [2 3]

	s2[1] = 999
	fmt.Println("После изменения S2 (s1 тоже изменился, т.к. общий массив):")
	fmt.Println(s1) // [1 2 999 4]
	fmt.Println(s2) // [2 999]

	s2 = append(s2, 5, 6, 7, 8, 9)
	fmt.Println("После добавления к S2 (произошла реалокация нового массива):")
	fmt.Println(s1)
	fmt.Println(s2)
}

// Strings демонстрирует итерацию по UTF-8 строке:
// `range s` декодирует руны (UTF-8 символы), возвращая индекс байта начала руны и код руны (int32).
func Strings() {
	s := "Привет!"

	for index, runeVal := range s {
		fmt.Printf("байт: %d, руна: %c (код: %d)\n", index, runeVal, runeVal)
	}
}

// CustomError — пользовательский тип ошибки для демонстрации классической ловушки `nil interface`.
type CustomError struct {
	Msg string
}

func (e *CustomError) Error() string {
	if e == nil {
		return "nil custom error"
	}
	return e.Msg
}

// FailReturningCustomError демонстрирует классическую ловушку `nil != nil` в Go:
// В Go интерфейс в рантайме представляется структурой `iface/eface` с двумя полями: (тип, значение).
// Интерфейс равен nil ТОЛЬКО когда оба поля (и тип, и значение) равны nil.
// Если интерфейсу присвоен типизированный указатель nil (например, `var err *CustomError = nil`),
// то интерфейс `error` получает тип `*CustomError` и значение `nil`.
// В результате проверка `err != nil` возвращает TRUE!
func FailReturningCustomError() error {
	var err *CustomError = nil
	return err // Ловушка: не-nil интерфейс, содержащий типизированный nil-указатель
}

// CorrectFail возвращает нетипизированный nil, когда ошибки нет.
func CorrectFail() error {
	var err *CustomError = nil
	if err == nil {
		return nil // Корректно: возвращает nil интерфейс
	}
	return err
}

// DeferWithNamedReturn демонстрирует поведение defer с именованными возвращаемыми значениями:
// Инструкция `return 10` присваивает значение переменной `res = 10`, затем вызывается `defer`,
// который удваивает значение `res = 20`, и только потом функция завершается.
func DeferWithNamedReturn() (res int) {
	defer func() {
		res *= 2
	}()
	return 10
}
