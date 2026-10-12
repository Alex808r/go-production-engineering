// Package main демонстрирует механизм утверждения типов (Type Assertion),
// приведение интерфейса к другому интерфейсу и влияние наборов методов (Method Sets).
package main

import (
	"fmt"
	"log"
)

// modern — интерфейс современной аудиоаппаратуры.
type modern interface {
	isModern() bool
}

// old — интерфейс устаревшей аппаратуры.
type old interface {
	isOld() bool
	isObsolete() bool
}

// cassete — аудиокассета.
type cassete struct {
	label string
}

// isModern имеет value receiver (c cassete).
// Значит, метод входит в набор методов как cassete, так и *cassete.
func (c cassete) isModern() bool {
	return false
}

// isOld имеет value receiver (c cassete).
// Входит в набор методов cassete и *cassete.
func (c cassete) isOld() bool {
	return true
}

// isObsolete имеет pointer receiver (c *cassete)!
// ВАЖНОЕ ПРАВИЛО GO (Method Sets):
// Методы с pointer receiver входят в набор методов ТОЛЬКО для указателя (*cassete),
// но НЕ входят в набор методов значения (cassete)!
// Поэтому:
//  - *cassete реализует интерфейс old (так как имеет и isOld, и isObsolete).
//  - cassete (по значению) НЕ реализует интерфейс old (у него нет isObsolete).
func (c *cassete) isObsolete() bool {
	return true
}

func main() {
	var m modern
	var c cassete = cassete{label: "Default"}

	// Присваиваем интерфейсу значение по значению
	m = c

	// m.label = "Ace Of Base" // ОШИБКА компиляции: через переменную интерфейса доступны ТОЛЬКО методы интерфейса!
	fmt.Println("m.isModern():", m.isModern())

	// 1. Утверждение типа к конкретной структуре (Type Assertion)
	// Синтаксис с проверкой: value, ok := iface.(ConcreteType)
	if casseteInstance, ok := m.(cassete); ok {
		casseteInstance.label = "Ace Of Base"
		log.Printf("Успешное утверждение типа: %+v", casseteInstance)
	}

	// 2. Утверждение типа к другому интерфейсу (Interface-to-Interface assertion)
	// Проверяем, реализует ли лежащий внутри объект другой интерфейс `old`.
	// Так как в `m` лежит копия по значению (cassete), а метод isObsolete() требует указателя,
	// проверка вернет false!
	if oldInstance, ok := m.(old); ok {
		fmt.Println("Old:", oldInstance.isOld(), "Obsolete:", oldInstance.isObsolete())
	} else {
		fmt.Println("cassete по значению НЕ реализует old (метод isObsolete требует указатель *cassete)")
	}

	// 3. А теперь проверим то же самое для указателя *cassete:
	var m2 modern = &c
	if oldInstance, ok := m2.(old); ok {
		fmt.Printf("Успех: *cassete реализует old! isOld=%v, isObsolete=%v\n",
			oldInstance.isOld(), oldInstance.isObsolete())
	}

	// 4. Опасность утверждения типа без проверки (comma-ok idiom):
	// В Go, если приводить неинтерфейсный тип к интерфейсу с методами, компилятор
	// проверяет возможность на этапе сборки: m.(int) выдаст ошибку компиляции "impossible type assertion".
	// А вот если тип скрыт за `any`, проверка происходит в runtime:
	var anyVal any = m
	// v := anyVal.(int) // Вызовет ПАНИКУ: panic: interface conversion: any is main.cassete, not int!
	// Поэтому всегда используйте безопасную форму с ok:
	if _, ok := anyVal.(int); !ok {
		fmt.Println("Безопасная проверка: anyVal не является int, паника предотвращена благодаря comma-ok.")
	}
}


