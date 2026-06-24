package operators_test

import (
	"reflect"
	"testing"
	"go-gap-engine/ga/operators"
)

func TestSinglePointCrossover(t *testing.T) {
	p1 := []int{1, 1, 1, 1}
	p2 := []int{2, 2, 2, 2}

	cx := operators.SinglePointCrossover[[]int, int]()
	o1, o2 := cx(p1, p2)

	if len(o1) != 4 || len(o2) != 4 {
		t.Errorf("Expected length 4, got %v and %v", len(o1), len(o2))
	}
}

func TestOrderCrossover(t *testing.T) {
	p1 := []int{1, 2, 3, 4, 5}
	p2 := []int{5, 4, 3, 2, 1}

	cx := operators.OrderCrossover[[]int, int]()
	o1, o2 := cx(p1, p2)

	// Quick check: should still contain all elements
	m1 := make(map[int]bool)
	for _, v := range o1 {
		m1[v] = true
	}
	if len(m1) != 5 {
		t.Errorf("Expected permutation to retain all 5 elements, got %v unique elements: %v", len(m1), o1)
	}
}

func TestBitFlip(t *testing.T) {
	ind := []bool{false, false, false}
	mut := operators.BitFlip[[]bool]()
	
	// rate 1.0 means all flip
	o := mut(ind, 1.0)
	if !reflect.DeepEqual(o, []bool{true, true, true}) {
		t.Errorf("Expected all bits to flip, got %v", o)
	}
}
