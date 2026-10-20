package geom

import (
	"math"
	"testing"
)

func TestPoint_Distance(t *testing.T) {
	tests := []struct {
		name     string
		p1       Point
		p2       Point
		expected float64
	}{
		{
			name:     "Египетский треугольник 3-4-5",
			p1:       Point{X: 1, Y: 1},
			p2:       Point{X: 4, Y: 5},
			expected: 5.0,
		},
		{
			name:     "Совпадающие точки (нулевое расстояние)",
			p1:       Point{X: 7.2, Y: -3.4},
			p2:       Point{X: 7.2, Y: -3.4},
			expected: 0.0,
		},
		{
			name:     "Отрицательные координаты",
			p1:       Point{X: -3, Y: -4},
			p2:       Point{X: 0, Y: 0},
			expected: 5.0,
		},
		{
			name:     "Симметричные точки по оси X",
			p1:       Point{X: -5, Y: 0},
			p2:       Point{X: 5, Y: 0},
			expected: 10.0,
		},
		{
			name:     "Симметричные точки по оси Y",
			p1:       Point{X: 0, Y: -6},
			p2:       Point{X: 0, Y: 6},
			expected: 12.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p1.Distance(tt.p2)
			if math.Abs(got-tt.expected) > 1e-9 {
				t.Errorf("Point.Distance() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGeom_BackwardCompatibility(t *testing.T) {
	g := Geom{X1: 1, Y1: 1, X2: 4, Y2: 5}

	if got := g.CalculateDistance(); math.Abs(got-5.0) > 1e-9 {
		t.Errorf("Geom.CalculateDistance() = %v, want 5.0", got)
	}

	if got := g.Distance(); math.Abs(got-5.0) > 1e-9 {
		t.Errorf("Geom.Distance() = %v, want 5.0", got)
	}
}

func BenchmarkPoint_Distance(b *testing.B) {
	p1 := Point{X: 1.5, Y: 2.5}
	p2 := Point{X: 4.5, Y: 6.5}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p1.Distance(p2)
	}
}
