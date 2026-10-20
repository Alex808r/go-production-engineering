// Пример реализации структуры данных "Двоичное дерево поиска".
// Пример для простоты приведён для дерева,
// содержащего в качетве значений целые числа.
// Можно по аналогии со стандартной библиотекой с помощью интерфейса
// обобщить на произвольные типы данных.
// Для этого потребуется контракт на функцию сравнения элементов.
package main

import "fmt"

// Двоичное дерево поиска (Binary Search Tree, BST).
// Свойство BST: для любого узла все значения в левом поддереве меньше значения узла,
// а все значения в правом поддереве — больше или равны.
// Временная сложность поиска и вставки:
//   - В среднем (сбалансированное дерево): O(log N)
//   - В худшем случае (вырождение в связный список): O(N)
type Tree struct {
	root *Element // Указатель на корень дерева (nil для пустого дерева)
}

// Element — узел двоичного дерева.
type Element struct {
	left, right *Element // Указатели на левого (меньшие) и правого (большие/равные) потомков
	Value       int      // Хранимое значение
}

// Insert вставляет новое значение в дерево за O(log N) в среднем.
func (t *Tree) Insert(x int) {
	e := &Element{Value: x}
	// Если дерево пустое, созданный элемент становится корнем
	if t.root == nil {
		t.root = e
		return
	}
	// Иначе спускаемся по дереву рекурсивно и ищем подходящее место
	insert(t.root, e)
}

// insert — вспомогательная рекурсивная функция для поиска листа и привязки нового узла.
func insert(node, new *Element) {
	// Если добавляемое значение меньше текущего — идем в левую ветку
	if new.Value < node.Value {
		if node.left == nil {
			node.left = new // нашли пустое место — привязываем левым потомком
			return
		}
		insert(node.left, new)
	}
	// Если добавляемое значение больше или равно — идем в правую ветку (дубликаты уходят вправо)
	if new.Value >= node.Value {
		if node.right == nil {
			node.right = new // нашли пустое место — привязываем правым потомком
			return
		}
		insert(node.right, new)
	}
}

// Search проверяет наличие значения x в дереве.
// Временная сложность: O(log N) в среднем, O(N) в худшем.
func (t *Tree) Search(x int) bool {
	return search(t.root, x)
}

// search — рекурсивный бинарный поиск по узлам дерева.
func search(el *Element, x int) bool {
	if el == nil {
		return false // дошли до листа и не нашли значение
	}
	if el.Value == x {
		return true // значение найдено!
	}
	// Если искомое больше текущего — отсекаем левую половину и идем вправо
	if el.Value < x {
		return search(el.right, x)
	}
	// Иначе отсекаем правую половину и идем влево
	return search(el.left, x)
}

// String реализует интерфейс fmt.Stringer для наглядного вывода дерева в консоль.
func (t Tree) String() string {
	return prettyPrint(t.root, 0)
}

// prettyPrint печатает дерево в консоль в виде горизонтальной диаграммы:
// Используется обратный симметричный обход (Reverse In-Order Traversal: Right -> Root -> Left):
// 1. Сначала рекурсивно печатается правое поддерево (большие числа окажутся вверху вывода).
// 2. Затем текущий узел с отступом табуляции, пропорциональным глубине уровня.
// 3. Затем рекурсивно левое поддерево (меньшие числа окажутся внизу вывода).
func prettyPrint(e *Element, spaces int) (res string) {
	if e == nil {
		return res
	}
	spaces++
	res += prettyPrint(e.right, spaces)
	for i := 0; i < spaces; i++ {
		res += "\t"
	}
	res += fmt.Sprintf("%d\n", e.Value)
	res += prettyPrint(e.left, spaces)
	return res
}

func initTree() *Tree {
	var t Tree
	t.Insert(10)
	t.Insert(15)
	t.Insert(15)
	t.Insert(15)
	t.Insert(15)
	t.Insert(25)
	t.Insert(30)
	t.Insert(35)
	t.Insert(5)
	t.Insert(20)
	t.Insert(1)
	t.Insert(2)
	t.Insert(6)
	return &t
}

func main() {
	t := initTree()
	fmt.Println("Визуализация дерева в консоли (повернуто на 90 градусов влево):")
	fmt.Println(t)

	// Демонстрация поиска
	fmt.Printf("Поиск 25: %v (ожидается true)\n", t.Search(25))
	fmt.Printf("Поиск 99: %v (ожидается false)\n", t.Search(99))
}

// 	Output:
// 					35
// 				30
// 			25
// 		20
// 			15
// 	10
//			6
// 		5
// 				2
// 			1
