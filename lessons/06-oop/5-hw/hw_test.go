package hw

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
			p1:       Point{X: 2.5, Y: 3.7},
			p2:       Point{X: 2.5, Y: 3.7},
			expected: 0.0,
		},
		{
			name:     "Отрицательные координаты (II и IV квадранты)",
			p1:       Point{X: -3, Y: -4},
			p2:       Point{X: 0, Y: 0},
			expected: 5.0,
		},
		{
			name:     "Симметричные точки через начало координат",
			p1:       Point{X: -3, Y: 0},
			p2:       Point{X: 3, Y: 0},
			expected: 6.0,
		},
		{
			name:     "Вертикальный отрезок",
			p1:       Point{X: 10, Y: -5},
			p2:       Point{X: 10, Y: 5},
			expected: 10.0,
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
	tests := []struct {
		name         string
		geom         Geom
		wantDistance float64
	}{
		{
			name:         "Базовый тест курса (#1)",
			geom:         Geom{X1: 1, Y1: 1, X2: 4, Y2: 5},
			wantDistance: 5.0,
		},
		{
			name:         "Отрицательные координаты через Geom",
			geom:         Geom{X1: -1, Y1: -1, X2: -4, Y2: -5},
			wantDistance: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.geom.CalculateDistance(); math.Abs(got-tt.wantDistance) > 1e-9 {
				t.Errorf("Geom.CalculateDistance() = %v, want %v", got, tt.wantDistance)
			}
			if got := tt.geom.Distance(); math.Abs(got-tt.wantDistance) > 1e-9 {
				t.Errorf("Geom.Distance() = %v, want %v", got, tt.wantDistance)
			}
		})
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
