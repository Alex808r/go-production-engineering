package isp_test

import (
	"testing"

	"go-core-4/lessons/17-system-design/SOLID/4-ISP"
)

type simpleDancer struct {
	danced bool
}

func (s *simpleDancer) Dance() {
	s.danced = true
}

func TestGoodDance_ISP(t *testing.T) {
	d := &simpleDancer{}
	// GoodDance принимает минимальный интерфейс Dancer (1 метод),
	// и simpleDancer не обязан реализовывать Sing() или PlayGuitar()
	isp.GoodDance(d)

	if !d.danced {
		t.Error("expected danced to be true")
	}
}
