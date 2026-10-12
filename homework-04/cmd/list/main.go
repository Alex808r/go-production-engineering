package main

import (
	"fmt"

	"go-core-4/homework-04/pkg/list"
)

func main() {
	fmt.Println("=== Демонстрация работы циклического двусвязного списка ===")

	l := list.New()
	l.Push(list.Elem{Val: 30})
	l.Push(list.Elem{Val: 20})
	l.Push(list.Elem{Val: 10})

	fmt.Printf("Исходный список: %s\n", l.String())

	l.Reverse()
	fmt.Printf("После разворота (Reverse): %s\n", l.String())

	l.Pop()
	fmt.Printf("После удаления первого (Pop): %s\n", l.String())

	// Цепочечный вызов Pop
	l.Pop().Pop()
	fmt.Printf("После l.Pop().Pop(): %q (пусто)\n", l.String())
}
