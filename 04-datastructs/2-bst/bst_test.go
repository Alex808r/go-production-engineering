package main

import (
	"strings"
	"testing"
)

func TestBST_InsertAndSearch(t *testing.T) {
	var tree Tree

	// Вставляем элементы
	values := []int{10, 5, 15, 3, 7, 12, 18}
	for _, v := range values {
		tree.Insert(v)
	}

	// Проверяем, что все добавленные элементы находятся
	for _, v := range values {
		if !tree.Search(v) {
			t.Errorf("Search(%d) = false; want true", v)
		}
	}

	// Проверяем отсутствующие элементы
	missing := []int{0, 1, 4, 8, 20, 100}
	for _, v := range missing {
		if tree.Search(v) {
			t.Errorf("Search(%d) = true; want false", v)
		}
	}
}

func TestBST_String(t *testing.T) {
	var tree Tree
	tree.Insert(10)
	tree.Insert(5)
	tree.Insert(15)

	out := tree.String()
	if !strings.Contains(out, "10") || !strings.Contains(out, "5") || !strings.Contains(out, "15") {
		t.Errorf("String() missing expected elements, got:\n%s", out)
	}
}
