package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteObject_LSP(t *testing.T) {
	buf := new(bytes.Buffer)

	s := String("Hello")
	b := Bool(true)
	a := Album{Title: "Master of Puppets", Year: 1986}

	// Все объекты подставляются вместо интерфейса Serializer без нарушения логики вызывающего кода
	err := WriteObject(buf, s, b, a)
	if err != nil {
		t.Fatalf("WriteObject failed: %v", err)
	}

	result := buf.String()
	if !strings.Contains(result, "Hello") {
		t.Errorf("expected 'Hello' in output, got: %s", result)
	}
	if !strings.Contains(result, "true") {
		t.Errorf("expected 'true' in output, got: %s", result)
	}
	if !strings.Contains(result, "Master of Puppets") {
		t.Errorf("expected album title in output, got: %s", result)
	}
}

func BenchmarkWriteObject(b *testing.B) {
	buf := new(bytes.Buffer)
	s := String("Benchmark")
	item := Album{Title: "Kind of Blue", Year: 1959}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_ = WriteObject(buf, s, item)
	}
}
