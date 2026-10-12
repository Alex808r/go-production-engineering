package main

import (
	"fmt"
)

// Course представляет учебный курс.
type Course struct {
	author string
}

// SetAuthor использует Pointer Receiver (*Course), поскольку метод мутирует (изменяет) состояние структуры.
// Если бы здесь был Value Receiver (Course), метод менял бы только локальную копию структуры.
func (c *Course) SetAuthor(name string) {
	c.author = name
}

// Author использует Value Receiver (Course), так как метод только читает данные (read-only)
// и не модифицирует поля структуры.
func (c Course) Author() string {
	return c.author
}

func main() {
	// new(Course) выделяет память под структуру Course, инициализирует её нулями и возвращает указатель (*Course).
	c := new(Course)
	c.SetAuthor("1")
	fmt.Println("c.Author():", c.Author()) // выведет: 1

	// course объявлен как значение (Course), а не указатель.
	var course Course

	// ВАЖНЫЙ СИНТАКСИЧЕСКИЙ САХАР Go:
	// Метод SetAuthor объявлен для указателя (*Course), но мы вызываем его на значении (course).
	// Поскольку переменная course адресуема (addressable), компилятор Go автоматически
	// преобразует вызов в: (&course).SetAuthor("2").
	course.SetAuthor("2")
	fmt.Println("course.Author():", course.Author()) // выведет: 2

	// ПРАВИЛО ВЫБОРА РЕСИВЕРА:
	// 1. Value Receiver: методы могут вызываться как на значениях (T), так и на указателях (*T).
	// 2. Pointer Receiver: методы могут вызываться на указателях (*T) и адресуемых значениях (&T).
	// 3. Если хотя бы один метод типа мутирует структуру — для консистентности принято делать
	//    все методы этого типа с указательными получателями (*T).
}

// ОТВЕТ НА ВОПРОС:
// Что будет выведено, если получатель SetAuthor сделать значением: func (c Course) SetAuthor(name string)?
// В строке course.SetAuthor("2") в метод передастся КОПИЯ структуры course.
// Поле c.author изменится у копии и сразу уничтожится при выходе из метода.
// Исходное поле course.author останется пустым, и course.Author() выведет пустую строку ""!
