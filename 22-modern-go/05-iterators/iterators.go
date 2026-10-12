// Package iterators демонстрирует революционную возможность Go 1.23+:
// Range-over-func (кастомные итераторы через iter.Seq и iter.Seq2).
package iterators

import (
	"cmp"
	"iter"
)

// TreeNode представляет узел бинарного дерева поиска.
type TreeNode[T cmp.Ordered] struct {
	Val   T
	Left  *TreeNode[T]
	Right *TreeNode[T]
}

// Insert добавляет элемент в бинарное дерево поиска.
func (n *TreeNode[T]) Insert(v T) *TreeNode[T] {
	if n == nil {
		return &TreeNode[T]{Val: v}
	}
	if v < n.Val {
		n.Left = n.Left.Insert(v)
	} else if v > n.Val {
		n.Right = n.Right.Insert(v)
	}
	return n
}

// All возвращает кастомный итератор iter.Seq[T] для симметричного обхода дерева (In-Order).
// В Go 1.23 функция iter.Seq[T] имеет сигнатуру func(yield func(T) bool).
// Если yield возвращает false (например, внутри цикла произошел break), обход немедленно прекращается!
func (n *TreeNode[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		var inOrder func(node *TreeNode[T]) bool
		inOrder = func(curr *TreeNode[T]) bool {
			if curr == nil {
				return true
			}
			// Обход левого поддерева
			if !inOrder(curr.Left) {
				return false
			}
			// Выдача текущего значения
			if !yield(curr.Val) {
				return false
			}
			// Обход правого поддерева
			return inOrder(curr.Right)
		}
		inOrder(n)
	}
}

// Slice возвращает классический срез для сравнительных бенчмарков.
func (n *TreeNode[T]) Slice() []T {
	var res []T
	var collect func(curr *TreeNode[T])
	collect = func(curr *TreeNode[T]) {
		if curr == nil {
			return
		}
		collect(curr.Left)
		res = append(res, curr.Val)
		collect(curr.Right)
	}
	collect(n)
	return res
}

// Filter — итератор-комбинатор высшего порядка: возвращает последовательность,
// отфильтрованную предикатом predicate.
func Filter[T any](seq iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range seq {
			if predicate(val) {
				if !yield(val) {
					return
				}
			}
		}
	}
}

// Pairs — демонстрация iter.Seq2[K, V] для итерации парами (ключ, значение).
func Pairs[K comparable, V any](m map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}
