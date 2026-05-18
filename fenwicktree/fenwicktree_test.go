package fenwicktree

import (
	"math"
	"testing"
)

func TestFenwickTreeEmpty(t *testing.T) {
	fw := NewFenwickTree[int64](0)
	if s := fw.Sum(0, 0); s != 0 {
		t.Errorf("expected 0, got %v", s)
	}
}

func TestFenwickTreeZero(t *testing.T) {
	fw := NewFenwickTree[int64](0)
	if s := fw.Sum(0, 0); s != 0 {
		t.Errorf("expected 0, got %v", s)
	}
}

func TestFenwickTreeOverFlowULL(t *testing.T) {
	fw := NewFenwickTree[uint64](10)
	for i := 0; i < 10; i++ {
		fw.Add(i, (1<<63)+uint64(i))
	}
	for i := 0; i <= 10; i++ {
		for j := i; j <= 10; j++ {
			var sum uint64 = 0
			for k := i; k < j; k++ {
				sum += uint64(k)
			}
			expected := sum
			if (j-i)%2 != 0 {
				expected = (1 << 63) + sum
			}
			if s := fw.Sum(i, j); s != expected {
				t.Errorf("expected %v, got %v", expected, s)
			}
		}
	}
}

func TestFenwickTreeNaiveTest(t *testing.T) {
	for n := 0; n <= 50; n++ {
		fw := NewFenwickTree[int64](n)
		for i := 0; i < n; i++ {
			fw.Add(i, int64(i*i))
		}
		for l := 0; l <= n; l++ {
			for r := l; r <= n; r++ {
				var sum int64 = 0
				for i := l; i < r; i++ {
					sum += int64(i * i)
				}
				if s := fw.Sum(l, r); s != sum {
					t.Errorf("expected %v, got %v", sum, s)
				}
			}
		}
	}
}

func TestFenwickTreeBound(t *testing.T) {
	fw := NewFenwickTree[int32](10)
	fw.Add(3, math.MaxInt32)
	fw.Add(5, math.MinInt32)
	if s := fw.Sum(0, 10); s != -1 {
		t.Errorf("expected -1, got %v", s)
	}
	if s := fw.Sum(3, 6); s != -1 {
		t.Errorf("expected -1, got %v", s)
	}
	if s := fw.Sum(3, 4); s != math.MaxInt32 {
		t.Errorf("expected MaxInt32, got %v", s)
	}
	if s := fw.Sum(4, 10); s != math.MinInt32 {
		t.Errorf("expected MinInt32, got %v", s)
	}
}

func TestFenwickTreeBoundLL(t *testing.T) {
	fw := NewFenwickTree[int64](10)
	fw.Add(3, math.MaxInt64)
	fw.Add(5, math.MinInt64)
	if s := fw.Sum(0, 10); s != -1 {
		t.Errorf("expected -1, got %v", s)
	}
	if s := fw.Sum(3, 6); s != -1 {
		t.Errorf("expected -1, got %v", s)
	}
	if s := fw.Sum(3, 4); s != math.MaxInt64 {
		t.Errorf("expected MaxInt64, got %v", s)
	}
	if s := fw.Sum(4, 10); s != math.MinInt64 {
		t.Errorf("expected MinInt64, got %v", s)
	}
}

func TestFenwickTreeOverFlow(t *testing.T) {
	fw := NewFenwickTree[int32](20)
	a := make([]int64, 20)
	for i := 0; i < 10; i++ {
		x := int32(math.MaxInt32)
		a[i] += int64(x)
		fw.Add(i, x)
	}
	for i := 10; i < 20; i++ {
		x := int32(math.MinInt32)
		a[i] += int64(x)
		fw.Add(i, x)
	}
	a[5] += 11111
	fw.Add(5, 11111)

	for l := 0; l <= 20; l++ {
		for r := l; r <= 20; r++ {
			var sum int64 = 0
			for i := l; i < r; i++ {
				sum += a[i]
			}
			dif := sum - int64(fw.Sum(l, r))
			if dif%(1<<32) != 0 {
				t.Errorf("dif %% (1<<32) should be 0, got %v", dif)
			}
		}
	}
}
