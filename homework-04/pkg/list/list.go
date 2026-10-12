// Package list реализует циклический двусвязный список со сторожевым элементом (sentinel node).
// Использование фиктивного корня позволяет избежать постоянных проверок на nil
// при добавлении и удалении граничных элементов.
package list

import (
	"fmt"
)

// List — циклический двусвязный список со сторожевым элементом root.
type List struct {
	root *Elem // Сторожевой элемент. Сам не хранит пользовательских данных.
}

// Elem — элемент (узел) двусвязного списка.
type Elem struct {
	Val        interface{} // Полезная нагрузка
	next, prev *Elem       // Указатели на следующий и предыдущий элементы
}

// New создаёт список и инициализирует сторожевой узел.
// В пустом кольцевом списке корень замыкается сам на себя: root.next = root и root.prev = root.
func New() *List {
	var l List
	l.root = &Elem{}
	l.root.next = l.root
	l.root.prev = l.root
	return &l
}

// Push вставляет элемент в НАЧАЛО списка (сразу после l.root) за O(1).
func (l *List) Push(e Elem) *Elem {
	// 1. Связываем новый элемент:
	//    - его prev указывает на корень;
	//    - его next указывает на бывший первый элемент (l.root.next).
	e.prev = l.root
	e.next = l.root.next

	// 2. Корень теперь указывает на новый элемент как на первый.
	l.root.next = &e

	// 3. Обновляем обратную связь у соседа справа:
	if e.next != l.root {
		e.next.prev = &e
	} else {
		// Если список был пуст, замыкаем l.root.prev на новый элемент (хвост списка).
		l.root.prev = &e
	}
	return &e
}

// String реализует интерфейс fmt.Stringer, обходя список от головы к хвосту за O(N).
func (l *List) String() string {
	el := l.root.next
	var s string
	// Идем по кольцу вперед, пока снова не упремся в сторожевой узел l.root
	for el != l.root {
		s += fmt.Sprintf("%v ", el.Val)
		el = el.next
	}
	if len(s) > 0 {
		s = s[:len(s)-1]
	}
	return s
}

// Pop удаляет первый элемент списка (голову) за O(1).
// Поддерживает цепочечные вызовы (fluent interface: l.Pop().Pop()).
func (l *List) Pop() *List {
	if l == nil || l.root == nil || l.root.next == l.root {
		return l
	}

	first := l.root.next
	l.root.next = first.next
	first.next.prev = l.root

	// Если список опустел, восстанавливаем указатель prev на корень
	if l.root.next == l.root {
		l.root.prev = l.root
	}

	return l
}

// Reverse разворачивает список на месте (in-place) за O(N) времени и O(1) памяти.
// Алгоритм меняет местами указатели next и prev у каждого узла.
func (l *List) Reverse() *List {
	if l == nil || l.root == nil || l.root.next == l.root {
		return l
	}

	curr := l.root.next     // начинаем с первого реального элемента (головы)
	var prev *Elem = l.root // предыдущим узлом для первого элемента изначально является root

	// Обходим все реальные узлы списка до тех пор, пока не вернемся в root
	for curr != l.root {
		next := curr.next // 1. Запоминаем оригинальный следующий узел
		curr.next = prev  // 2. Разворачиваем стрелку next назад (на prev)
		curr.prev = next  // 3. Разворачиваем стрелку prev вперед (на оригинальный next)
		prev = curr       // 4. Сдвигаем prev вперед
		curr = next       // 5. Переходим к следующему узлу
	}

	// Корректируем сторожевой узел root после полного разворота:
	oldHead := l.root.next
	l.root.next = prev
	l.root.prev = oldHead

	return l
}
