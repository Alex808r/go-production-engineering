package main

import (
	"fmt"
	"iter"
	"math/big"
	"strconv"
)

// 1. РЕКУРСИЯ С МЕМОИЗАЦИЕЙ (Top-Down динамическое программирование)
// Временная сложность: O(n) благодаря кэшу (каждое поддерево вычисляется ровно один раз).
// Без кэширования наивная рекурсия имела бы экспоненциальную сложность O(2^n).
// Пространственная сложность: O(n) — память под хэш-таблицу cache плюс глубина стека вызовов n.
func fibo(n int, cache map[int]int) int {
	// Проверка наличия уже вычисленного результата в кэше
	if val, ok := cache[n]; ok {
		return val
	}
	// Базовая часть рекурсии
	if n < 2 {
		return 1
	}
	// Рекурсивный шаг с записью в кэш
	num := fibo(n-1, cache) + fibo(n-2, cache)
	cache[n] = num
	return num
}

// 2. ИТЕРАТИВНЫЙ ЦИКЛ (Bottom-Up динамическое программирование)
// Временная сложность: O(n) — один линейный проход.
// Пространственная сложность: O(1) — константная память, используются только 2 переменные.
// Идеальный выбор для продакшена: полностью отсутствует риск переполнения стека (Stack Overflow).
func fiboIterative(n int) int {
	if n < 2 {
		return 1
	}
	a, b := 1, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b // сдвиг окна из двух последних значений
	}
	return b
}

// 3. ГЕНЕРАТОР ЧЕРЕЗ ЗАМЫКАНИЕ (Closure из Tour of Go)
// Функция захватывает переменные a и b в куче (Escape to Heap).
// При каждом последующем вызове возвращает очередное число последовательности.
// Временная сложность: O(1) за каждый вызов. Пространственная сложность: O(1).
func fiboClosure() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

// 4. СОВРЕМЕННЫЙ ИТЕРАТОР Go 1.23+ (iter.Seq)
// Ленивая бесконечная последовательность, совместимая со стандартным циклом for ... range.
// Вычисление выполняется по требованию (on-demand), без выделения памяти под срез.
func fiboSeq() iter.Seq[int] {
	return func(yield func(int) bool) {
		a, b := 1, 1
		for {
			if !yield(a) { // yield вернёт false, если вызывающий цикл выполнил break
				return
			}
			a, b = b, a+b
		}
	}
}

// 5. ВЫЧИСЛЕНИЕ ДЛЯ БОЛЬШИХ ЧИСЕЛ (math/big)
// Тип int (64 бита) переполняется уже на F(92)/F(93):
// результат тихо оборачивается в отрицательное число.
// math.Big позволяет работать с числами произвольной разрядности.
func fiboBig(n int) *big.Int {
	if n < 2 {
		return big.NewInt(1)
	}
	a := big.NewInt(1)
	b := big.NewInt(1)
	next := new(big.Int)

	for i := 2; i <= n; i++ {
		next.Add(a, b)
		a.Set(b)
		b.Set(next)
	}
	return b
}

func main() {
	fmt.Println("=== Способы вычисления чисел Фибоначчи в Go ===")

	// Демонстрация 1: Итератор Go 1.23+
	fmt.Print("1. Ленивый итератор Go 1.23+ (первые 10 чисел): ")
	count := 0
	for v := range fiboSeq() {
		fmt.Printf("%d ", v)
		count++
		if count >= 10 {
			break
		}
	}
	fmt.Println()

	// Демонстрация 2: Замыкание (Closure)
	fmt.Print("2. Генератор через замыкание (первые 10 чисел):  ")
	gen := fiboClosure()
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", gen())
	}
	fmt.Println()

	// Демонстрация 3: Вычисление больших чисел без переполнения int64
	fmt.Printf("3. Большое число Фибоначчи F(100) через math/big:\n   %s\n\n", fiboBig(100).String())

	// Демонстрация 4: Интерактивный расчет и сравнение подходов
	cache := make(map[int]int)
	for {
		fmt.Print("Введите число n (или 'exit' для выхода): ")
		var in string
		_, err := fmt.Scanln(&in)
		if err != nil || in == "exit" || in == "q" {
			fmt.Println("Выход.")
			break
		}
		n, err := strconv.Atoi(in)
		if err != nil || n < 0 {
			fmt.Println("Пожалуйста, введите неотрицательное целое число.")
			continue
		}

		fmt.Printf("   fibo (мемоизация):  F(%d) = %d\n", n, fibo(n, cache))
		fmt.Printf("   fiboIterative:      F(%d) = %d\n", n, fiboIterative(n))
		if n >= 92 {
			fmt.Printf("   fiboBig (math/big): F(%d) = %s (для предотвращения переполнения int64)\n", n, fiboBig(n).String())
		}
	}
}
