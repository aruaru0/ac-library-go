package twosat

import (
	"math/rand"
	"testing"
)

func TestTwoSATEmpty(t *testing.T) {
	ts := NewTwoSAT(0)
	if !ts.Satisfiable() {
		t.Errorf("Empty 2SAT should be satisfiable")
	}
	if len(ts.Answer()) != 0 {
		t.Errorf("Empty 2SAT answer length should be 0")
	}
}

func TestTwoSATOne(t *testing.T) {
	{
		ts := NewTwoSAT(1)
		ts.AddClause(0, true, 0, true)
		ts.AddClause(0, false, 0, false)
		if ts.Satisfiable() {
			t.Errorf("Should not be satisfiable")
		}
	}
	{
		ts := NewTwoSAT(1)
		ts.AddClause(0, true, 0, true)
		if !ts.Satisfiable() {
			t.Errorf("Should be satisfiable")
		}
		if ts.Answer()[0] != true {
			t.Errorf("Expected true")
		}
	}
	{
		ts := NewTwoSAT(1)
		ts.AddClause(0, false, 0, false)
		if !ts.Satisfiable() {
			t.Errorf("Should be satisfiable")
		}
		if ts.Answer()[0] != false {
			t.Errorf("Expected false")
		}
	}
}

func TestTwoSATStressOK(t *testing.T) {
	for phase := 0; phase < 1000; phase++ {
		n := rand.Intn(20) + 1
		m := rand.Intn(100) + 1
		expect := make([]bool, n)
		for i := 0; i < n; i++ {
			expect[i] = rand.Intn(2) == 1
		}
		ts := NewTwoSAT(n)
		xs := make([]int, m)
		ys := make([]int, m)
		types := make([]int, m)
		for i := 0; i < m; i++ {
			x := rand.Intn(n)
			y := rand.Intn(n)
			typ := rand.Intn(3)
			xs[i] = x
			ys[i] = y
			types[i] = typ
			if typ == 0 {
				ts.AddClause(x, expect[x], y, expect[y])
			} else if typ == 1 {
				ts.AddClause(x, !expect[x], y, expect[y])
			} else {
				ts.AddClause(x, expect[x], y, !expect[y])
			}
		}
		if !ts.Satisfiable() {
			t.Errorf("Should be satisfiable")
		}
		actual := ts.Answer()
		for i := 0; i < m; i++ {
			x, y, typ := xs[i], ys[i], types[i]
			if typ == 0 {
				if !(actual[x] == expect[x] || actual[y] == expect[y]) {
					t.Errorf("Clause 0 failed")
				}
			} else if typ == 1 {
				if !(actual[x] != expect[x] || actual[y] == expect[y]) {
					t.Errorf("Clause 1 failed")
				}
			} else {
				if !(actual[x] == expect[x] || actual[y] != expect[y]) {
					t.Errorf("Clause 2 failed")
				}
			}
		}
	}
}
