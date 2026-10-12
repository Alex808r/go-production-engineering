package main

import (
	"fmt"
)

// Пример иллюстрирует композицию (встраивание) структур и затенение методов (Method Shadowing).
// В Go нет наследования классов (extends), вместо него используется встраивание (embedding).

// computer — внешняя структура, анонимно встраивающая processor.
type computer struct {
	model string
	processor // Встраивание: поля и методы processor "продвигаются" на уровень computer.
}

// processor — встраиваемая структура (компонент).
type processor struct {
	model string
	cores int
}

// cpuinfo у типа computer переопределяет (затеняет) метод cpuinfo встроенного процессора.
func (c *computer) cpuinfo() string {
	return "computer cpuinfo()"
}

// cpuinfo у типа processor возвращает подробные характеристики процессора.
func (p *processor) cpuinfo() string {
	return fmt.Sprintf("%s, %d ядер", p.model, p.cores)
}

// String реализует интерфейс fmt.Stringer для красивого форматирования.
func (c *computer) String() string {
	// Обратите внимание: c.processor.cpuinfo() явно вызывает метод встроенной структуры,
	// минуя затенение внешним методом c.cpuinfo().
	return fmt.Sprintf("Компьютер: %s\nПроцессор: %s", c.model, c.processor.cpuinfo())
}

func main() {
	comp := &computer{
		model: "Компьютер игровой",
		processor: processor{
			model: "Байкал-М",
			cores: 8,
		},
	}

	// 1. Затенение методов (Shadowing):
	// comp.cpuinfo() вызывает метод типа computer, так как он находится на верхнем уровне.
	fmt.Println("1. comp.cpuinfo() [внешний]:", comp.cpuinfo())

	// 2. Прямой доступ к встроенному методу:
	// Чтобы вызвать метод встроенной структуры, обращаемся через имя встроенного типа.
	fmt.Println("2. comp.processor.cpuinfo() [вложенный]:", comp.processor.cpuinfo())

	// 3. Продвижение полей (Field Promotion):
	// К полю cores можно обратиться напрямую через comp.cores (продвинуто),
	// а к модели процессора — через comp.processor.model (так как comp.model перекрывает имя).
	fmt.Printf("3. Продвинутое поле comp.cores: %d\n\n", comp.cores)

	// 4. Форматированный вывод всего объекта через интерфейс fmt.Stringer:
	fmt.Println(comp)
}
