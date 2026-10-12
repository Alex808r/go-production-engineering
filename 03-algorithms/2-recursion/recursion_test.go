package main

import "testing"

func Test_pow(t *testing.T) {
	tests := []struct {
		name string
		x    int
		y    int
		want int
	}{
		{"base 0", 2, 0, 1},
		{"positive", 2, 3, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pow(tt.x, tt.y)
			if got != tt.want {
				t.Errorf("pow() = %v, want %v", got, tt.want)
			}
		})
	}
}
