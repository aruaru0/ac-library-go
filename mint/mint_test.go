package mint

import (
	"reflect"
	"testing"
)

func TestModIntBasic(t *testing.T) {
	m := NewModInt(13)

	// Add
	if m.Add(5, 10) != 2 {
		t.Errorf("Add(5, 10) should be 2, got %d", m.Add(5, 10))
	}
	if m.Add(5, -2) != 3 {
		t.Errorf("Add(5, -2) should be 3, got %d", m.Add(5, -2))
	}

	// Sub
	if m.Sub(5, 10) != 8 {
		t.Errorf("Sub(5, 10) should be 8, got %d", m.Sub(5, 10))
	}

	// Mul
	if m.Mul(4, 5) != 7 {
		t.Errorf("Mul(4, 5) should be 7, got %d", m.Mul(4, 5))
	}

	// ModInv
	if m.ModInv(4) != 10 { // 4 * 10 = 40 = 3 * 13 + 1
		t.Errorf("ModInv(4) should be 10, got %d", m.ModInv(4))
	}

	// Div
	if m.Div(5, 4) != m.Mul(5, 10) {
		t.Errorf("Div(5, 4) should be %d, got %d", m.Mul(5, 10), m.Div(5, 4))
	}

	// Pow
	if m.Pow(3, 4) != 3 { // 3^4 = 81 = 6 * 13 + 3
		t.Errorf("Pow(3, 4) should be 3, got %d", m.Pow(3, 4))
	}
}

func TestModIntCombination(t *testing.T) {
	m := NewModInt(998244353)

	// nCr basic
	if m.NCr(5, 2) != 10 {
		t.Errorf("NCr(5, 2) should be 10, got %d", m.NCr(5, 2))
	}
	if m.NCr(5, 0) != 1 {
		t.Errorf("NCr(5, 0) should be 1, got %d", m.NCr(5, 0))
	}
	if m.NCr(5, 5) != 1 {
		t.Errorf("NCr(5, 5) should be 1, got %d", m.NCr(5, 5))
	}
	if m.NCr(5, 6) != 0 {
		t.Errorf("NCr(5, 6) should be 0, got %d", m.NCr(5, 6))
	}
	if m.NCr(5, -1) != 0 {
		t.Errorf("NCr(5, -1) should be 0, got %d", m.NCr(5, -1))
	}

	// Large nCr test (ensure no TLE)
	m.InitComb(1000)
	if m.NCr(1000, 500) == 0 {
		t.Errorf("NCr(1000, 500) should not be 0")
	}
}

func TestModIntMatrix(t *testing.T) {
	m := NewModInt(1000)

	A := [][]int{
		{1, 1},
		{1, 0},
	}

	// Fibonacci using matrix power (A^5)
	// A^0 = I
	// A^1 = [1, 1; 1, 0]
	// A^2 = [2, 1; 1, 1]
	// A^3 = [3, 2; 2, 1]
	// A^4 = [5, 3; 3, 2]
	// A^5 = [8, 5; 5, 3]
	expect := [][]int{
		{8, 5},
		{5, 3},
	}
	actual := m.PowModMatrix(A, 5)

	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("PowModMatrix(A, 5) failed, expected %v, got %v", expect, actual)
	}
}
