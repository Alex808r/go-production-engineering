package list

import (
	"testing"
)

func TestList_Push(t *testing.T) {
	l := New()
	l.Push(Elem{Val: 2})
	l.Push(Elem{Val: 1})

	got := l.String()
	want := "1 2"
	if got != want {
		t.Fatalf("получили %s, ожидалось %s", got, want)
	}
}

func TestList_Pop(t *testing.T) {
	l := New()
	l.Push(Elem{Val: 3})
	l.Push(Elem{Val: 2})
	l.Push(Elem{Val: 1})

	got := l.Pop().String()
	want := "2 3"
	if got != want {
		t.Fatalf("получили %s, ожидалось %s", got, want)
	}

	got = l.Pop().Pop().String()
	want = ""
	if got != want {
		t.Fatalf("получили %s, ожидалось %s", got, want)
	}

	// Pop на пустом списке не должен приводить к панике
	got = l.Pop().String()
	if got != "" {
		t.Fatalf("Pop на пустом списке: получили %q, ожидалось \"\"", got)
	}
}

func TestList_Reverse(t *testing.T) {
	l := New()
	l.Push(Elem{Val: 3})
	l.Push(Elem{Val: 2})
	l.Push(Elem{Val: 1})

	got := l.Reverse().String()
	want := "3 2 1"
	if got != want {
		t.Fatalf("получили %s, ожидалось %s", got, want)
	}

	// Повторный разворот возвращает исходный порядок
	got = l.Reverse().String()
	want = "1 2 3"
	if got != want {
		t.Fatalf("повторный reverse: получили %s, ожидалось %s", got, want)
	}
}
