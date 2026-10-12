package users

import (
	"fmt"
	"io"
)

// ============================================================================
// Задача №1. Полиморфизм через интерфейс с методом-геттером
// ============================================================================

// Ager — контракт для типов, предоставляющих информацию о возрасте.
type Ager interface {
	Age() int
}

// Employee — сотрудник с методом-геттером Age().
type Employee struct {
	Name string
	age  int
}

// NewEmployee создает новый экземпляр Employee.
func NewEmployee(name string, age int) Employee {
	return Employee{Name: name, age: age}
}

// Age возвращает возраст сотрудника.
func (e Employee) Age() int {
	return e.age
}

// Customer — клиент с методом-геттером Age().
type Customer struct {
	Name string
	age  int
}

// NewCustomer создает новый экземпляр Customer.
func NewCustomer(name string, age int) Customer {
	return Customer{Name: name, age: age}
}

// Age возвращает возраст клиента.
func (c Customer) Age() int {
	return c.age
}

// MaxAge принимает произвольное количество объектов любого типа,
// реализующего контракт Ager, и возвращает возраст самого старшего человека.
// Если список пуст, возвращает 0.
func MaxAge(people ...Ager) int {
	var maxAge int
	for _, p := range people {
		if p == nil {
			continue
		}
		if a := p.Age(); a > maxAge {
			maxAge = a
		}
	}
	return maxAge
}

// ============================================================================
// Задача №2. Работа с пустым интерфейсом any и Type Switch
// ============================================================================

// EmployeeSimple — сотрудник без методов.
type EmployeeSimple struct {
	Name string
	Age  int
}

// CustomerSimple — клиент без методов.
type CustomerSimple struct {
	Name string
	Age  int
}

// OldestPerson принимает произвольное количество аргументов любого типа (any),
// извлекает возраст структур EmployeeSimple и CustomerSimple с помощью переключателя типов (type switch)
// и возвращает сам объект наиболее старшего человека.
// Если список пуст или подходящих объектов нет, возвращает nil.
func OldestPerson(people ...any) any {
	var oldest any
	maxAge := -1

	for _, p := range people {
		switch v := p.(type) {
		case EmployeeSimple:
			if v.Age > maxAge {
				maxAge = v.Age
				oldest = v
			}
		case CustomerSimple:
			if v.Age > maxAge {
				maxAge = v.Age
				oldest = v
			}
		case *EmployeeSimple:
			if v != nil && v.Age > maxAge {
				maxAge = v.Age
				oldest = v
			}
		case *CustomerSimple:
			if v != nil && v.Age > maxAge {
				maxAge = v.Age
				oldest = v
			}
		default:
			// Игнорируем типы, не являющиеся EmployeeSimple или CustomerSimple
		}
	}

	return oldest
}

// ============================================================================
// Задача №3. Фильтрация типов и потоковый вывод в io.Writer
// ============================================================================

// WriteOnlyStrings принимает io.Writer и произвольное количество аргументов любого типа.
// Функция фильтрует аргументы и выводит в io.Writer ТОЛЬКО строки.
// Возвращает общее количество записанных байт и ошибку записи (если возникла).
func WriteOnlyStrings(w io.Writer, args ...any) (int, error) {
	if w == nil {
		return 0, fmt.Errorf("writer cannot be nil")
	}

	var totalWritten int

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			n, err := io.WriteString(w, v)
			totalWritten += n
			if err != nil {
				return totalWritten, fmt.Errorf("write string failed: %w", err)
			}
		case *string:
			if v != nil {
				n, err := io.WriteString(w, *v)
				totalWritten += n
				if err != nil {
					return totalWritten, fmt.Errorf("write string pointer failed: %w", err)
				}
			}
		default:
			// Все остальные типы (числа, булевы, структуры и т.д.) пропускаются
		}
	}

	return totalWritten, nil
}
