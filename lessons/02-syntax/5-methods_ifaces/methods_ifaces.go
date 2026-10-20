package main

import "fmt"

// WalkerTalker демонстрирует композицию (встраивание) интерфейсов:
// он объединяет методы Walk(), Talk() и все методы вложенного интерфейса Eater (Eat()).
type WalkerTalker interface {
	Walk()        // контракт интерфейса
	Talk() string // контракт интерфейса
	Eater         // встраивание интерфейса Eater
}

// MyIface — пустой интерфейс. Поскольку у него 0 методов,
// абсолютно любой тип в Go удовлетворяет этому контракту (аналог interface{} / any).
type MyIface interface {
}

type Eater interface {
	Eat() string
}

// Guy — пользовательский тип (структура).
// Встраивание интерфейса Eater наделяет Guy методом Eat() string на уровне типов,
// благодаря чему Guy формально удовлетворяет интерфейсу Eater без явного написания метода.
type Guy struct {
	Eater
}

// Метод с pointer receiver (*Guy).
// Примечание: так как ресивер указательный (*Guy), контракт интерфейса реализует
// именно указатель *Guy, а не само значение Guy (правило Method Sets).
func (g *Guy) Walk() {
	fmt.Println("I walk!")
}

// Метод с pointer receiver (*Guy).
func (g *Guy) Talk() string {
	return "I talk!"
}

func main() {
	// Переменная интерфейсного типа инициализирована
	// литералом конкретного типа.
	var g Guy

	// 1. Полиморфизм и утиная типизация (Duck Typing):
	//    Интерфейс WalkerTalker требует 3 метода: Walk(), Talk() и Eat().
	//    У типа *Guy есть все 3 метода (Walk и Talk объявлены явно, Eat получен через встраивание).
	//    Поэтому компилятор считает, что *Guy реализует WalkerTalker, без явного слова "implements".
	//
	// 2. Почему именно адрес &g, а не значение g:
	//    Методы объявлены с ресивером-указателем: func (g *Guy) Walk().
	//    В Go значение типа Guy НЕ наследует методы типа *Guy.
	//    Поэтому "var wt WalkerTalker = g" не скомпилируется (Guy does not implement WalkerTalker).
	//
	// 3. Динамическая диспетчеризация:
	//    Переменная wt под капотом хранит пару [конкретный тип: *Guy, адрес данных: &g].
	//    При вызове wt.Walk() и wt.Talk() рантайм находит нужные методы структуры *Guy и вызывает их.
	var wt WalkerTalker = &g
	wt.Walk()          // вызывает (*Guy).Walk() -> печатает "I walk!"
	println(wt.Talk()) // вызывает (*Guy).Talk() -> возвращает "I talk!", println выводит в консоль

	// ВНИМАНИЕ: вызов wt.Eat() или g.Eat() вызовет panic (nil pointer dereference),
	// так как встроенное поле-интерфейс Eater у переменной g не проинициализировано (равно nil).
	// println(wt.Eat())

	// Демонстрация пустого интерфейса: можно присвоить значение любого типа
	var mi MyIface
	i := 10
	s := "str"
	mi = i
	mi = s
	_ = mi

	// g.Eat() // Также вызовет panic, если поле Eater не заполнено конкретной реализацией

	// Стандартный пустой интерфейс interface{} (в Go 1.18+ есть алиас any)
	var in interface{}
	in = "20"
	in = "abc"
	in = []int{1, 2, 3}
	_ = in
}

// Интерфейс с одним методом
type printer interface {
	print()
}

// Пользовательский тип на основе базового типа string.
// В Go методы можно определять для ЛЮБЫХ именованных типов, а не только для структур.
type myStr string

// Метод с pointer receiver (*myStr) позволяет мутировать значение по указателю.
func (s *myStr) print() {
	*s = "DEF"
	fmt.Println(*s)
}
