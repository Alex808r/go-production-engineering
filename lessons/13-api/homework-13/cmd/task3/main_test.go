package main

import (
	"testing"
	"unsafe"

	"github.com/DmitriyVTitov/size"
)

func TestTask3_Sizes(t *testing.T) {
	ex := example{
		a: []int{1, 2, 3},
		b: true,
		d: "1234",
	}

	ex1 := example{
		a: []int{1, 2, 3},
		b: true,
		d: "1234",
		c: 100,
	}

	if unsafe.Sizeof(ex) != 48 {
		t.Errorf("unsafe.Sizeof(ex) = %d, want 48", unsafe.Sizeof(ex))
	}
	if unsafe.Sizeof(ex1) != 48 {
		t.Errorf("unsafe.Sizeof(ex1) = %d, want 48", unsafe.Sizeof(ex1))
	}

	if size.Of(ex) != 76 {
		t.Errorf("size.Of(ex) = %d, want 76", size.Of(ex))
	}
	if size.Of(ex1) != 76 {
		t.Errorf("size.Of(ex1) = %d, want 76", size.Of(ex1))
	}
}
