package basic

import (
	"fmt"
	"reflect"
)

func Range() {
	slice := []int{10, 20, 30, 40}
	arr := [...]int{10, 20, 30}
	// arr2 := [2]int{1,2}
	str := "АБВ"

	for i, v := range slice {
		fmt.Println(i, v)
	}
	fmt.Printf("\n\n")

	for i, v := range arr {
		fmt.Println(i, v)
	}
	fmt.Printf("\n\n")

	for i, v := range str {
		fmt.Println(i, v, string(v))
	}
	fmt.Printf("\n\n")
}

func Vars() {
	var i int
	var j int8 = 10
	k := 20
	var l = 10
	println(i, j, k, l)
}

func Types() {
	type circle struct {
		x, y int
		r    int
	}
	c := circle{
		x: 10,
		y: 10,
		r: 20,
	}
	fmt.Printf("%+v\n", c)

	type myString string
	var s string = "ABC"
	var s1 myString = myString(s)
	fmt.Printf("Тип 1: %v\tТип2: %v\t\n", reflect.TypeOf(s), reflect.TypeOf(s1))
}

func Pointers() {
	var s string = "ABC"
	var pointer *string = &s
	fmt.Println(*pointer, pointer)
	fmt.Printf("Значение: %v\tУказатель: %v\tСсылка: %v\t\n", s, pointer, &s)
}

// Scopes демонстрирует блочные области видимости (block scope) и затенение переменных (variable shadowing).
// Вложенный блок создает новую переменную x, затеняя внешнюю, при этом значение внешней переменной не изменяется.
func Scopes() {
	x := 1
	{
		// Внутренний блок: объявляется новая локальная переменная x, затеняющая внешнюю
		x := 2
		fmt.Printf("x = %v\n", x) // выведет: x = 2
	}
	// Внешний блок: исходная переменная x осталась неизменной
	fmt.Printf("x = %v\n", x) // выведет: x = 1
}
