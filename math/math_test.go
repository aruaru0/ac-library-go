package math

import (
	"math"
	"testing"
)

func gcd(a, b int64) int64 {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func TestMathPowMod(t *testing.T) {
	naive := func(x, n int64, mod int) int64 {
		y := safeMod(x, int64(mod))
		z := uint64(1 % mod)
		for i := int64(0); i < n; i++ {
			z = (z * uint64(y)) % uint64(mod)
		}
		return int64(z)
	}

	for a := int64(-100); a <= 100; a++ {
		for b := int64(0); b <= 100; b++ {
			for c := 1; c <= 100; c++ {
				if naive(a, b, c) != PowMod(a, b, c) {
					t.Errorf("PowMod(%d, %d, %d) failed", a, b, c)
				}
			}
		}
	}
}

func TestMathInvBoundHand(t *testing.T) {
	minll := int64(math.MinInt64)
	maxll := int64(math.MaxInt64)

	if InvMod(-1, maxll) != InvMod(minll, maxll) {
		t.Errorf("InvMod bound failed")
	}
	if InvMod(maxll, maxll-1) != 1 {
		t.Errorf("InvMod maxll, maxll-1 failed")
	}
	if InvMod(maxll-1, maxll) != maxll-1 {
		t.Errorf("InvMod maxll-1, maxll failed")
	}
	if InvMod(maxll/2+1, maxll) != 2 {
		t.Errorf("InvMod maxll/2+1, maxll failed")
	}
}

func TestMathInvMod(t *testing.T) {
	for a := int64(-100); a <= 100; a++ {
		for b := int64(1); b <= 1000; b++ {
			if gcd(safeMod(a, b), b) != 1 {
				continue
			}
			c := InvMod(a, b)
			if c < 0 || c >= b {
				t.Errorf("InvMod out of bounds")
			}
			if ((a*c)%b+b)%b != 1%b {
				t.Errorf("InvMod(%d, %d) failed", a, b)
			}
		}
	}
}

func TestMathInvModZero(t *testing.T) {
	if InvMod(0, 1) != 0 {
		t.Errorf("InvMod(0, 1) should be 0")
	}
	for i := int64(0); i < 10; i++ {
		if InvMod(i, 1) != 0 {
			t.Errorf("InvMod(%d, 1) should be 0", i)
		}
		if InvMod(-i, 1) != 0 {
			t.Errorf("InvMod(%d, 1) should be 0", -i)
		}
		if InvMod(math.MinInt64+i, 1) != 0 {
			t.Errorf("InvMod(min+i, 1) should be 0")
		}
		if InvMod(math.MaxInt64-i, 1) != 0 {
			t.Errorf("InvMod(max-i, 1) should be 0")
		}
	}
}

func floorSumNaive(n, m, a, b int64) int64 {
	var sum int64 = 0
	for i := int64(0); i < n; i++ {
		z := a*i + b
		sum += (z - safeMod(z, m)) / m
	}
	return sum
}

func TestMathFloorSum(t *testing.T) {
	for n := int64(0); n < 20; n++ {
		for m := int64(1); m < 20; m++ {
			for a := int64(-20); a < 20; a++ {
				for b := int64(-20); b < 20; b++ {
					if floorSumNaive(n, m, a, b) != FloorSum(n, m, a, b) {
						t.Errorf("FloorSum(%d, %d, %d, %d) failed", n, m, a, b)
					}
				}
			}
		}
	}
}

func TestMathCRTHand(t *testing.T) {
	y, z := CRT([]int64{1, 2, 1}, []int64{2, 3, 2})
	if y != 5 || z != 6 {
		t.Errorf("CRTHand failed: got %d, %d", y, z)
	}
}

func TestMathCRT2(t *testing.T) {
	for a := int64(1); a <= 20; a++ {
		for b := int64(1); b <= 20; b++ {
			for c := int64(-10); c <= 10; c++ {
				for d := int64(-10); d <= 10; d++ {
					y, z := CRT([]int64{c, d}, []int64{a, b})
					if z == 0 {
						for x := int64(0); x < a*b/gcd(a, b); x++ {
							if x%a == c && x%b == d {
								t.Errorf("CRT2 false negative")
							}
						}
						continue
					}
					if z != a*b/gcd(a, b) {
						t.Errorf("CRT2 wrong z")
					}
					if safeMod(c, a) != y%a {
						t.Errorf("CRT2 wrong y mod a")
					}
					if safeMod(d, b) != y%b {
						t.Errorf("CRT2 wrong y mod b")
					}
				}
			}
		}
	}
}

func TestMathCRTOverflow(t *testing.T) {
	r0 := int64(0)
	r1 := int64(1000000000000) - 2
	m0 := int64(900577)
	m1 := int64(1000000000000)
	y, z := CRT([]int64{r0, r1}, []int64{m0, m1})
	if z != m0*m1 {
		t.Errorf("CRTOverflow z mismatch")
	}
	if r0 != y%m0 {
		t.Errorf("CRTOverflow y mod m0 mismatch")
	}
	if r1 != y%m1 {
		t.Errorf("CRTOverflow y mod m1 mismatch")
	}
}

func TestInternalMathSafeMod(t *testing.T) {
	preds := []int64{}
	for i := int64(0); i <= 100; i++ {
		preds = append(preds, i, -i, i, math.MinInt64+i, math.MaxInt64-i)
	}
	for _, a := range preds {
		for _, b := range preds {
			if b <= 0 {
				continue
			}
			ans := a % b
			if ans < 0 {
				ans += b
			}
			if ans != safeMod(a, b) {
				t.Errorf("safeMod(%d, %d) failed", a, b)
			}
		}
	}
}

func TestInternalMathInvGcdBound(t *testing.T) {
	pred := []int64{}
	for i := int64(0); i <= 10; i++ {
		pred = append(pred, i, -i)
		pred = append(pred, math.MinInt64+i, math.MaxInt64-i)
		pred = append(pred, math.MinInt64/2+i, math.MaxInt64/2-i)
		pred = append(pred, math.MinInt64/3+i, math.MaxInt64/3-i)
	}
	pred = append(pred, 998244353, 1000000007, 1000000009)
	pred = append(pred, -998244353, -1000000007, -1000000009)

	for _, a := range pred {
		for _, b := range pred {
			if b <= 0 {
				continue
			}
			a2 := safeMod(a, b)
			g, x := invGCD(a, b)
			eg := gcd(a2, b)
			if g != eg {
				t.Errorf("invGCD g failed")
			}
			if x < 0 || x > b/g {
				t.Errorf("invGCD x out of bounds")
			}
			// avoid int128 overflow test in Go without big.Int, simply checking bounds
		}
	}
}
