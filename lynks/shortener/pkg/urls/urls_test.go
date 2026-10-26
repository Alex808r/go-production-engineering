package urls

import (
	"testing"
)

func TestShorten(t *testing.T) {
	s1 := Shorten("https://github.com/Alex808r/go-production-engineering")
	if len(s1) != 5 {
		t.Fatalf("expected length 5, got %d (%s)", len(s1), s1)
	}

	s2 := Shorten("https://github.com/Alex808r/go-production-engineering")
	if len(s2) != 5 {
		t.Fatalf("expected length 5, got %d (%s)", len(s2), s2)
	}

	if Shorten("") != "" {
		t.Fatalf("expected empty string for empty input")
	}
}

func BenchmarkShorten(b *testing.B) {
	url := "https://github.com/Alex808r/go-production-engineering"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Shorten(url)
	}
}
